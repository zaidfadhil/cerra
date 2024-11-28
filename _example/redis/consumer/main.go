package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zaidfadhil/cerra"
	"github.com/zaidfadhil/cerra/redisq"
)

func handleTask(ctx context.Context, t *cerra.Task) error {
	fmt.Println("get", t)
	return nil
}

func main() {
	redisQueue := redisq.New(redisq.Options{
		Address:  "localhost:6379",
		Password: "redis",
		Stream:   "cerra",
		Group:    "cerra",
		Consumer: "cerra",
	})
	queue := cerra.NewQueue(redisQueue, 2)

	queue.AddHandler(handleTask)

	done := shutdown(func(grace bool, done chan<- struct{}) {
		queue.Close()
		done <- struct{}{}
	})

	queue.Start()

	<-done
}

func shutdown(stop func(done chan<- struct{})) <-chan struct{} {
	done := make(chan struct{})
	s := make(chan os.Signal, 1)
	signal.Notify(s,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		<-s
		stop(done)
		close(done)
	}()
	return done
}
