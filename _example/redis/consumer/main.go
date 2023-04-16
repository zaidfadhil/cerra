package main

import (
	"context"
	"fmt"

	"github.com/zaidfadhil/goatq"
	"github.com/zaidfadhil/goatq/redisq"
)

func handleTask(ctx context.Context, t *goatq.Task) error {
	fmt.Println("get", string(t.Payload))
	return nil
}

func main() {
	m := goatq.NewManager()

	opts := redisq.RedisOptions{
		Address:  "localhost:6379",
		Password: "redis",
		Queue:    "goatq",
		Group:    "goatq",
		Consumer: "goatq",
	}

	queue := goatq.NewQueue(redisq.NewRedisBackend(opts))

	queue.AddHandler(handleTask)

	queue.Start()

	m.OnShutdown(func() error {
		queue.Close()
		return nil
	})

	<-m.Done()
}
