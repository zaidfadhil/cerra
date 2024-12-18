package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zaidfadhil/cerra"
)

type Model struct {
	Title string
	Num   int
}

func newTask(title string, num int) (*cerra.Task, error) {
	payload, err := json.Marshal(&Model{Title: title, Num: num})
	if err != nil {
		return nil, err
	}
	fmt.Println("set:", payload)
	return cerra.NewTask(payload), nil
}

func handleTask(ctx context.Context, t *cerra.Task) error {
	var model Model
	if err := json.Unmarshal(t.Payload, &model); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}
	fmt.Println("get:", model.Num)
	return nil
}

func main() {
	queue := cerra.NewQueue(cerra.NewInMemoryBackend(), 2)

	for i := 0; i < 100; i++ {
		task, err := newTask("test-queue", i)
		if err != nil {
			fmt.Println(err)
		}

		err = queue.Enqueue(task)
		if err != nil {
			fmt.Println(err)
		}
	}

	queue.AddHandler(handleTask)

	done := shutdown(func(done chan<- struct{}) {
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
