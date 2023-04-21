package main

import (
	"context"
	"fmt"

	"github.com/zaidfadhil/goatq"
	"github.com/zaidfadhil/goatq/rabbitmq"
)

func handleTask(ctx context.Context, t *goatq.Task) error {
	fmt.Println("get", string(t.Payload))
	return nil
}

func main() {
	m := goatq.NewManager()

	rabbitQueue := rabbitmq.New(rabbitmq.Options{
		Address:      "amqp://user:pass@localhost:5672/",
		Queue:        "goatq",
		ExchangeName: "goatq-exchange",
		ExchangeType: "direct",
		RoutingKey:   "goatq-key",
	})
	queue := goatq.NewQueue(rabbitQueue)

	queue.AddHandler(handleTask)

	queue.Start()

	m.OnShutdown(func() error {
		queue.Close()
		return nil
	})

	<-m.Done()
}
