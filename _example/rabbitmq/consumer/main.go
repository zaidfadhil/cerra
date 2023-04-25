package main

import (
	"context"
	"fmt"

	"github.com/zaidfadhil/cerra"
	"github.com/zaidfadhil/cerra/rabbitmq"
)

func handleTask(ctx context.Context, t *cerra.Task) error {
	fmt.Println("get", string(t.Payload))
	return nil
}

func main() {
	m := cerra.NewManager()

	rabbitQueue := rabbitmq.New(rabbitmq.Options{
		Address:      "amqp://user:pass@localhost:5672/",
		Queue:        "cerra",
		ExchangeName: "cerra-exchange",
		ExchangeType: "direct",
		RoutingKey:   "cerra-key",
	})
	queue := cerra.NewQueue(rabbitQueue, 2)

	queue.AddHandler(handleTask)

	queue.Start()

	m.OnShutdown(func() error {
		queue.Close()
		return nil
	})

	<-m.Done()
}
