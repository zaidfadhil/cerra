package main

import (
	"fmt"

	"github.com/zaidfadhil/goatq"
	redisq "github.com/zaidfadhil/goatq/redis"
)

func main() {
	opts := redisq.RedisOptions{
		Address:  "localhost:6379",
		Password: "redis",
		Queue:    "goatq",
		Group:    "goatq",
		Consumer: "goatq",
	}

	queue := goatq.NewQueue(redisq.NewRedisBackend(opts))

	for i := 0; i < 100000; i++ {
		task := goatq.NewTask("test", []byte(fmt.Sprint(i)))
		err := queue.Enqueue(task)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(i)
	}
}
