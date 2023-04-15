package redisq

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zaidfadhil/goatq"
)

type RedisOptions struct {
	Client *redis.Client

	ConnectionString string
	Address          string
	Cluster          bool
	DB               int
	Password         string
	Queue            string
	Group            string
	Consumer         string

	blockTime time.Duration
}

var _ goatq.Backend = (*redisBackend)(nil)

type redisBackend struct {
	options RedisOptions

	rdb   redis.Cmdable
	tasks chan redis.XMessage

	stop      chan struct{}
	exit      chan struct{}
	startSync sync.Once
	stopSync  sync.Once
}

func NewRedisBackend(options RedisOptions) *redisBackend {
	b := &redisBackend{
		stop:    make(chan struct{}),
		exit:    make(chan struct{}),
		tasks:   make(chan redis.XMessage),
		options: options,
	}
	b.options.blockTime = 60 * time.Second

	if b.options.Client != nil {
		b.rdb = b.options.Client
	} else if b.options.ConnectionString != "" {
		options, err := redis.ParseURL(b.options.ConnectionString)
		if err != nil {
			log.Fatalf("error parsing redis connection string %v", err)
		}
		b.rdb = redis.NewClient(options)
	} else if b.options.Address != "" {
		if b.options.Cluster {
			b.rdb = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:    strings.Split(b.options.Address, ","),
				Password: b.options.Password,
			})
		} else {
			options := &redis.Options{
				Addr:     b.options.Address,
				Password: b.options.Password,
				DB:       b.options.DB,
			}
			b.rdb = redis.NewClient(options)
		}
	}

	_, err := b.rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("error connecting to redis %v", err)
	}

	b.stop = make(chan struct{})
	b.exit = make(chan struct{})
	b.tasks = make(chan redis.XMessage)
	b.options.blockTime = 60 * time.Second
	return b
}

func (b *redisBackend) Enqueue(task *goatq.Task) error {
	ctx := context.Background()

	err := b.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: b.options.Queue,
		Values: task.ToMap(),
	}).Err()

	return err
}

func (b *redisBackend) Dequeue() (*goatq.Task, error) {
	err := b.consumer()
	if err != nil {
		return nil, goatq.ErrInActiveQueue
	}

	task, ok := <-b.tasks
	if !ok {
		return nil, goatq.ErrInActiveQueue
	}

	return &goatq.Task{
		Name:    task.Values["name"].(string),
		Payload: []byte(task.Values["payload"].(string)),
	}, nil
}

func (b *redisBackend) Close() error {
	b.stopSync.Do(func() {
		close(b.stop)

		select {
		case <-b.exit:
		case <-time.After(200 * time.Millisecond):
		}

		switch rediscon := b.rdb.(type) {
		case *redis.Client:
			rediscon.Close()
		case *redis.ClusterClient:
			rediscon.Close()
		}
		close(b.tasks)
	})
	return nil
}

func (b *redisBackend) consumer() (err error) {
	b.startSync.Do(func() {
		ctx := context.Background()
		err := b.rdb.XGroupCreateMkStream(
			ctx,
			b.options.Queue,
			b.options.Group,
			"0",
		).Err()
		if err != nil {
			log.Printf("error creating stream %v", err)
		}

		go b.fetch()
	})
	return err
}

func (b *redisBackend) fetch() {

	for {
		select {
		case <-b.stop:
			return
		default:
		}

		ctx := context.Background()
		data, err := b.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    b.options.Group,
			Consumer: b.options.Consumer,
			Streams:  []string{b.options.Queue, ">"},
			Count:    1,
			Block:    b.options.blockTime,
		}).Result()
		if err != nil {
			log.Printf("error reading redis stream %v", err)
			continue
		}

		for _, result := range data {
			for _, message := range result.Messages {

				select {
				case b.tasks <- message:

					err := b.rdb.XAck(
						ctx,
						b.options.Queue,
						b.options.Group,
						message.ID,
					).Err()
					if err != nil {
						log.Printf("error message ack %v", err)
					}

					err = b.rdb.XDel(
						ctx,
						b.options.Queue,
						message.ID,
					).Err()
					if err != nil {
						log.Printf("error when deleting the message %v", err)
					}
				case <-b.stop:
					err := b.rdb.XAdd(ctx, &redis.XAddArgs{
						Stream: b.options.Queue,
						Values: message.Values,
					}).Err()
					if err != nil {
						log.Printf("error requeue the message %v", err)
					}
					close(b.exit)
					return
				}
			}
		}
	}
}
