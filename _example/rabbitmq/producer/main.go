package main

import (
	"fmt"
	"time"

	"github.com/zaidfadhil/cerra"
	"github.com/zaidfadhil/cerra/rabbitmq"
)

func main() {

	rabbitQueue := rabbitmq.New(rabbitmq.Options{
		Address:      "amqp://user:pass@localhost:5672/",
		Queue:        "cerra",
		ExchangeName: "cerra-exchange",
		ExchangeType: "direct",
		RoutingKey:   "cerra-key",
	})
	queue := cerra.NewQueue(rabbitQueue)

	for i := 0; i < 100; i++ {
		task := cerra.NewTask("test", []byte(fmt.Sprint(i)))
		err := queue.Enqueue(task)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(i)
	}

	time.Sleep(1 * time.Second)
}
