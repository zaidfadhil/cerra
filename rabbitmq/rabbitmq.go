package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zaidfadhil/goatq"
)

type Options struct {
	Address      string
	Queue        string
	ExchangeName string
	ExchangeType string
	RoutingKey   string
}

var _ goatq.Backend = (*rabbiMQBackend)(nil)

type rabbiMQBackend struct {
	options Options

	connection *amqp.Connection
	channel    *amqp.Channel
	tasks      <-chan amqp.Delivery

	stop      chan struct{}
	startSync sync.Once
	stopSync  sync.Once
}

func New(options Options) *rabbiMQBackend {
	b := &rabbiMQBackend{
		tasks:   make(chan amqp.Delivery),
		stop:    make(chan struct{}),
		options: defaultOptions(options),
	}
	var err error

	b.connection, err = amqp.Dial(options.Address)
	if err != nil {
		log.Fatalf("amqp dial error %v", err)
	}

	b.channel, err = b.connection.Channel()
	if err != nil {
		log.Fatalf("amqp connection error: %v", err)
	}

	err = b.channel.ExchangeDeclare(
		options.ExchangeName,
		options.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("amqp exchange declare error: %v", err)
	}

	return b
}

func (b *rabbiMQBackend) Enqueue(task *goatq.Task) error {
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
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            encodedTask,
			DeliveryMode:    amqp.Transient,
			Priority:        0,
		})
}

func (b *rabbiMQBackend) Dequeue() (*goatq.Task, error) {
	err := b.consumer()
	if err != nil {
		return nil, goatq.ErrInActiveQueue
	}

	task, ok := <-b.tasks
	if !ok {
		return nil, goatq.ErrInActiveQueue
	}

	var data goatq.Task
	_ = json.Unmarshal(task.Body, &data)
	_ = task.Ack(false)
	return &data, nil
}

func (b *rabbiMQBackend) Close() (err error) {
	b.stopSync.Do(func() {
		close(b.stop)
		if err = b.channel.Cancel(b.options.Queue, true); err != nil {
			log.Printf("rabbitmq channel close error: %v", err)
		}

		if err = b.connection.Close(); err != nil {
			log.Printf("rabbitmq connection close error: %v", err)
		}
	})
	return err
}

func (b *rabbiMQBackend) consumer() (err error) {
	b.startSync.Do(func() {
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
			return
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
			return
		}

		b.tasks, err = b.channel.Consume(
			q.Name,
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

func defaultOptions(opts Options) Options {
	if opts.Address == "" {
		opts.Address = "amqp://user:pass@localhost:5672"
	}
	if opts.Queue == "" {
		opts.Queue = "goatq-queue"
	}
	if opts.ExchangeName == "" {
		opts.ExchangeName = "goatq-exchange"
	}
	if opts.ExchangeType == "" {
		opts.ExchangeType = "direct"
	}
	if opts.RoutingKey == "" {
		opts.RoutingKey = "goatq-key"
	}
	return opts
}
