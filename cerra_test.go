package cerra_test

import (
	"context"
	"testing"
	"time"

	"github.com/zaidfadhil/cerra"
)

func TestEnqueue(t *testing.T) {
	queue := cerra.NewQueue(cerra.NewInMemoryBackend(), 1)

	task := &cerra.Task{
		ID:      "test_task",
		Payload: []byte("test_payload"),
	}

	err := queue.Enqueue(task)
	if err != nil {
		t.Errorf("enqueu error: %v", err)
	}
}

func TestAddHandler(t *testing.T) {
	queue := cerra.NewQueue(cerra.NewInMemoryBackend(), 1)

	var dequeuedTask *cerra.Task

	queue.AddHandler(func(ctx context.Context, t *cerra.Task) error {
		dequeuedTask = t
		return nil
	})

	queue.Start()
	defer queue.Close()

	task := &cerra.Task{
		ID:      "test_task",
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

	if dequeuedTask.ID != task.ID {
		t.Error("dequeue task name != queued task name")
	}

	if string(dequeuedTask.Payload) != string(task.Payload) {
		t.Error("dequeue task payload != queued task payload")
	}
}
