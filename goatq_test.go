package goatq_test

import (
	"context"
	"testing"
	"time"

	"github.com/zaidfadhil/goatq"
)

func TestEnqueue(t *testing.T) {
	queue := goatq.NewQueue(goatq.NewInMemoryBackend())

	task := &goatq.Task{
		Name:    "test_task",
		Payload: []byte("test_payload"),
	}

	err := queue.Enqueue(task)
	if err != nil {
		t.Errorf("enqueu error: %v", err)
	}
}

func TestAddHandler(t *testing.T) {
	queue := goatq.NewQueue(goatq.NewInMemoryBackend())

	var dequeuedTask *goatq.Task

	queue.AddHandler(func(ctx context.Context, t *goatq.Task) error {
		dequeuedTask = t
		return nil
	})

	queue.Start()
	defer queue.Close()

	task := &goatq.Task{
		Name:    "test_task",
		Payload: []byte("test_payload"),
	}

	err := queue.Enqueue(task)
	if err != nil {
		t.Errorf("enqueu error: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	if dequeuedTask == nil {
		t.Error("handler was not called")
	}

	if dequeuedTask.Name != task.Name {
		t.Error("dequeue task name != queued task name")
	}

	if string(dequeuedTask.Payload) != string(task.Payload) {
		t.Error("dequeue task payload != queued task payload")
	}
}
