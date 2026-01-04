package amqp

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zaidfadhil/cerra"
)

type Options struct {
	Client *amqp.Connection

	Address      string
	Queue        string
	ExchangeName string
	ExchangeType string
	RoutingKey   string
}

var _ cerra.Backend = (*amqpBackend)(nil)

type amqpBackend struct {
	options Options

	connection *amqp.Connection
	channel    *amqp.Channel
	tasks      <-chan amqp.Delivery

	stop      chan struct{}
	startSync sync.Once
	stopSync  sync.Once
}

func New(options Options) *amqpBackend {
	b := &amqpBackend{
		tasks:   make(chan amqp.Delivery),
		stop:    make(chan struct{}),
		options: defaultOptions(options),
	}
	var err error

	if b.options.Client != nil {
		b.connection = b.options.Client
	} else {
		b.connection, err = amqp.Dial(b.options.Address)
		if err != nil {
			log.Fatalf("amqp dial error %v", err)
		}
	}

	b.channel, err = b.connection.Channel()
	if err != nil {
		log.Fatalf("amqp connection error: %v", err)
	}

	err = b.channel.ExchangeDeclare(
		b.options.ExchangeName,
		b.options.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("amqp exchange declare error: %v", err)
	}

	_, err = b.bind()
	if err != nil {
		log.Fatalf("amqp bind error: %v", err)
	}

	return b
}

func (b *amqpBackend) Enqueue(task *cerra.Task) error {
	encodedTask, err := task.Encode()
	if err != nil {
		return err
	}
	return b.channel.PublishWithContext(
		context.Background(),
		b.options.ExchangeName,
		b.options.RoutingKey,
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{"task_id": task.ID},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            encodedTask,
			DeliveryMode:    amqp.Transient,
			Priority:        0,
		})
}

func (b *amqpBackend) Dequeue() (*cerra.Task, error) {
	err := b.consumer()
	if err != nil {
		return nil, cerra.ErrInActiveQueue
	}

	times := 0
loop:
	for {
		select {
		case task, ok := <-b.tasks:
			if !ok {
				return nil, cerra.ErrInActiveQueue
			}
			var data cerra.Task
			_ = json.Unmarshal(task.Body, &data)
			_ = task.Ack(false)
			return &data, nil
		case <-time.After(500 * time.Millisecond):
			if times == 5 {
				break loop
			}
			times += 1
		}
	}

	return nil, cerra.ErrEmtpyQueue
}

func (b *amqpBackend) Close() (err error) {
	b.stopSync.Do(func() {
		close(b.stop)
		if err = b.channel.Cancel(b.options.Queue, true); err != nil {
			log.Printf("amqp channel close error: %v", err)
		}

		if err = b.connection.Close(); err != nil {
			log.Printf("amqp connection close error: %v", err)
		}
	})
	return err
}

func (b *amqpBackend) consumer() (err error) {
	b.startSync.Do(func() {
		qName, err := b.bind()
		if err != nil {
			log.Println(err)
			return
		}

		b.tasks, err = b.channel.Consume(
			qName,
			b.options.Queue,
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Printf("amqp consumer error: %v", err)
			return
		}

	})
	return err
}

func (b *amqpBackend) bind() (string, error) {
	q, err := b.channel.QueueDeclare(
		b.options.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
		return "", err
	}

	err = b.channel.QueueBind(
		q.Name,
		b.options.RoutingKey,
		b.options.ExchangeName,
		false,
		nil,
	)
	if err != nil {
		log.Printf("exchange bind error: %v", err)
		return "", err
	}

	return q.Name, nil
}

func defaultOptions(opts Options) Options {
	if opts.Queue == "" {
		opts.Queue = "cerra-queue"
	}
	if opts.ExchangeName == "" {
		opts.ExchangeName = "cerra-exchange"
	}
	if opts.ExchangeType == "" {
		opts.ExchangeType = "direct"
	}
	if opts.RoutingKey == "" {
		opts.RoutingKey = "cerra-key"
	}
	return opts
}
