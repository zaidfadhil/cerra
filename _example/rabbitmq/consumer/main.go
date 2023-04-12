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

	opts := rabbitmq.RabbiMQOptions{
		Address:      "amqp://user:password@localhost:5672/",
		Queue:        "goatq",
		ExchangeName: "goatq-exchange",
		ExchangeType: "direct",
		RoutingKey:   "goatq-key",
	}

	queue := goatq.NewQueue(rabbitmq.NewRabbitMQBackend(opts))

	queue.AddHandler(handleTask)

	queue.Start()

	<-m.Done()
}
