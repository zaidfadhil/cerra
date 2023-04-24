package main

import (
	"context"
	"encoding/json"
	"fmt"

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
	return cerra.NewTask("yo:queue", payload), nil
}

func handleTask(ctx context.Context, t *cerra.Task) error {
	var model Model
	if err := json.Unmarshal(t.Payload, &model); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}
	fmt.Println("get-1:", model.Num)
	return nil
}

func main() {
	m := cerra.NewManager()

	queue := cerra.NewQueue(cerra.NewInMemoryBackend())

	for i := 0; i < 1000; i++ {
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

	queue.Start()

	m.OnShutdown(func() error {
		queue.Close()
		return nil
	})

	<-m.Done()
}
