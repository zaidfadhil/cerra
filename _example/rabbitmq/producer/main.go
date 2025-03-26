package main

import (
	"fmt"
	"time"

	"github.com/zaidfadhil/cerra"
	"github.com/zaidfadhil/cerra/amqp"
)

func main() {
	rabbitQueue := amqp.New(amqp.Options{
		Address:      "amqp://user:pass@localhost:5672/",
		Queue:        "cerra",
		ExchangeName: "cerra-exchange",
		ExchangeType: "direct",
		RoutingKey:   "cerra-key",
	})
	queue := cerra.NewQueue(rabbitQueue, 0)

	for i := 0; i < 100; i++ {
		task := cerra.NewTask([]byte(fmt.Sprint(i)))
		err := queue.Enqueue(task)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(i)
	}

	time.Sleep(1 * time.Second)
}
