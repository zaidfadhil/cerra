package main

import (
	"fmt"

	"github.com/zaidfadhil/goatq"
	"github.com/zaidfadhil/goatq/rabbitmq"
)

func main() {
	opts := rabbitmq.RabbiMQOptions{
		Address:      "amqp://user:password@localhost:5672",
		Queue:        "goatq",
		ExchangeName: "goatq-exchange",
		ExchangeType: "direct",
		RoutingKey:   "goatq-key",
	}

	queue := goatq.NewQueue(rabbitmq.NewRabbitMQBackend(opts))

	for i := 0; i < 1000; i++ {
		task := goatq.NewTask("test", []byte(fmt.Sprint(i)))
		err := queue.Enqueue(task)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(i)
	}
}
