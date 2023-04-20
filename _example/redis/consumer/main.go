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

	redisQueue := redisq.New(redisq.Options{
		Address:  "localhost:6379",
		Password: "redis",
		Stream:   "goatq",
		Group:    "goatq",
		Consumer: "goatq",
	})
	queue := goatq.NewQueue(redisQueue)

	queue.AddHandler(handleTask)

	queue.Start()

	m.OnShutdown(func() error {
		queue.Close()
		return nil
	})

	<-m.Done()
}
