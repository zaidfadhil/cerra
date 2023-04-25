package main

import (
	"context"
	"fmt"

	"github.com/zaidfadhil/cerra"
	"github.com/zaidfadhil/cerra/redisq"
)

func handleTask(ctx context.Context, t *cerra.Task) error {
	fmt.Println("get", string(t.Payload))
	return nil
}

func main() {
	m := cerra.NewManager()

	redisQueue := redisq.New(redisq.Options{
		Address:  "localhost:6379",
		Password: "redis",
		Stream:   "cerra",
		Group:    "cerra",
		Consumer: "cerra",
	})
	queue := cerra.NewQueue(redisQueue, 2)

	queue.AddHandler(handleTask)

	queue.Start()

	m.OnShutdown(func() error {
		queue.Close()
		return nil
	})

	<-m.Done()
}
