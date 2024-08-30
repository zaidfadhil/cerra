package main

import (
	"fmt"
	"time"

	"github.com/zaidfadhil/cerra"
	"github.com/zaidfadhil/cerra/redisq"
)

func main() {

	redisQueue := redisq.New(redisq.Options{
		Address:  "localhost:6379",
		Password: "redis",
		Stream:   "cerra",
		Group:    "cerra",
		Consumer: "cerra",
	})
	queue := cerra.NewQueue(redisQueue, 0)

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
