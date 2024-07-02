package redisq

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zaidfadhil/cerra"
)

type Options struct {
	Client *redis.Client

	ConnectionString string
	Address          string
	Cluster          bool
	DB               int
	Password         string
	Stream           string
	Group            string
	Consumer         string

	blockTime time.Duration
}

var _ cerra.Backend = (*redisBackend)(nil)

type redisBackend struct {
	options Options

	rdb   redis.Cmdable
	tasks chan redis.XMessage

	stop      chan struct{}
	exit      chan struct{}
	startSync sync.Once
	stopSync  sync.Once
}

func New(options Options) *redisBackend {
	b := &redisBackend{
		stop:    make(chan struct{}),
		exit:    make(chan struct{}),
		tasks:   make(chan redis.XMessage),
		options: defaultOptions(options),
	}

	if b.options.Client != nil {
		b.rdb = b.options.Client
	} else if b.options.ConnectionString != "" {
		options, err := redis.ParseURL(b.options.ConnectionString)
		if err != nil {
			log.Fatalf("error parsing redis connection string: %v", err)
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
		log.Fatalf("error connecting to redis: %v", err)
	}

	return b
}

func (b *redisBackend) Enqueue(task *cerra.Task) error {
	ctx := context.Background()

	err := b.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: b.options.Stream,
		Values: task.ToMap(),
	}).Err()

	return err
}

func (b *redisBackend) Dequeue() (*cerra.Task, error) {
	err := b.consumer()
	if err != nil {
		return nil, cerra.ErrInActiveQueue
	}

	task, ok := <-b.tasks
	if !ok {
		return nil, cerra.ErrInActiveQueue
	}

	return &cerra.Task{
		ID:      task.Values["id"].(string),
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
		err := b.createGroup(ctx)
		if err != nil {
			if !strings.HasPrefix(err.Error(), "BUSYGROUP") {
				log.Printf("error creating stream: %v", err)
			}
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
			Streams:  []string{b.options.Stream, ">"},
			Count:    1,
			Block:    b.options.blockTime,
		}).Result()
		if err != nil {
			if strings.HasPrefix(err.Error(), "NOGROUP") {
				b.createGroup(ctx)
			}
			log.Printf("error reading redis stream: %v", err)
			continue
		}

		for _, result := range data {
			for _, message := range result.Messages {
				select {
				case b.tasks <- message:
					b.ack(ctx, message)
				case <-b.stop:
					log.Printf("requeue %v, %v", message.ID, message.Values)
					err := b.rdb.XAdd(ctx, &redis.XAddArgs{
						Stream: b.options.Stream,
						Values: message.Values,
					}).Err()
					if err != nil {
						log.Printf("error requeue the message: %v", err)
					}
					close(b.exit)
					return
				}
			}
		}
	}
}

func (b *redisBackend) createGroup(ctx context.Context) error {
	return b.rdb.XGroupCreateMkStream(
		ctx,
		b.options.Stream,
		b.options.Group,
		"0",
	).Err()
}

func (b *redisBackend) ack(ctx context.Context, m redis.XMessage) {
	err := b.rdb.XAck(
		ctx,
		b.options.Stream,
		b.options.Group,
		m.ID,
	).Err()
	if err != nil {
		log.Printf("error message ack: %v", err)
	}

	err = b.rdb.XDel(
		ctx,
		b.options.Stream,
		m.ID,
	).Err()
	if err != nil {
		log.Printf("error when deleting the message: %v", err)
	}
}

func defaultOptions(opts Options) Options {
	if opts.Address == "" {
		opts.Address = "localhost:6379"
	}
	if opts.Stream == "" {
		opts.Stream = "cerra-stream"
	}
	if opts.Group == "" {
		opts.Group = "cerra-group"
	}
	if opts.Consumer == "" {
		opts.Consumer = "cerra-consumer"
	}
	return opts
}
