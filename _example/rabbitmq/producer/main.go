package main

import (
	"fmt"
	"time"

	"github.com/zaidfadhil/goatq"
	"github.com/zaidfadhil/goatq/rabbitmq"
)

func main() {

	rabbitQueue := rabbitmq.New(rabbitmq.Options{
		Address:      "amqp://user:pass@localhost:5672/",
		Queue:        "goatq",
		ExchangeName: "goatq-exchange",
		ExchangeType: "direct",
		RoutingKey:   "goatq-key",
	})
	queue := goatq.NewQueue(rabbitQueue)

	for i := 0; i < 100; i++ {
		task := goatq.NewTask("test", []byte(fmt.Sprint(i)))
		err := queue.Enqueue(task)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(i)
	}

	time.Sleep(1 * time.Second)
}
