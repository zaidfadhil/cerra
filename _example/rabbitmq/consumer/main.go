package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zaidfadhil/cerra"
	"github.com/zaidfadhil/cerra/amqp"
)

func handleTask(ctx context.Context, t *cerra.Task) error {
	fmt.Println("get", string(t.Payload))
	return nil
}

func main() {
	rabbitQueue := amqp.New(amqp.Options{
		Address:      "amqp://user:pass@localhost:5672/",
		Queue:        "cerra",
		ExchangeName: "cerra-exchange",
		ExchangeType: "direct",
		RoutingKey:   "cerra-key",
	})
	queue := cerra.NewQueue(rabbitQueue, 2)

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
