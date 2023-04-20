package main

import (
	"fmt"
	"time"

	"github.com/zaidfadhil/goatq"
	"github.com/zaidfadhil/goatq/redisq"
)

func main() {

	redisQueue := redisq.New(redisq.Options{
		Address:  "localhost:6379",
		Password: "redis",
		Stream:   "goatq",
		Group:    "goatq",
		Consumer: "goatq",
	})
	queue := goatq.NewQueue(redisQueue)

	for i := 0; i < 100000; i++ {
		task := goatq.NewTask("test", []byte(fmt.Sprint(i)))
		err := queue.Enqueue(task)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(i)
	}

	time.Sleep(1 * time.Second)
}
