package redisq_test

import (
	"context"
	"testing"
	"time"

	"github.com/zaidfadhil/cerra"
	"github.com/zaidfadhil/cerra/redisq"
)

func TestRedisEnqueue(t *testing.T) {
	backend := redisq.New(redisq.Options{
		Address: "localhost:6379",
	})
	queue := cerra.NewQueue(backend, 1)
	defer queue.Close()

	task := &cerra.Task{
		ID:      "test_task",
		Payload: []byte("test_payload"),
	}

	err := queue.Enqueue(task)
	if err != nil {
		t.Errorf("rabbitmq enqueu error: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
}

func TestRedisDequeue(t *testing.T) {
	backend := redisq.New(redisq.Options{
		Address: "localhost:6379",
	})
	queue := cerra.NewQueue(backend, 1)
	queue.Start()

	task := &cerra.Task{
		ID:      "test_task",
		Payload: []byte("test_payload"),
	}

	err := queue.Enqueue(task)
	if err != nil {
		t.Errorf("rabbitmq enqueue error: %v", err)
	}

	var dequeuedTask *cerra.Task
	queue.AddHandler(func(ctx context.Context, t *cerra.Task) error {
		dequeuedTask = t
		return nil
	})

	time.Sleep(100 * time.Millisecond)

	if dequeuedTask == nil {
		t.Error("handler was not called")
	}

	if dequeuedTask.ID != task.ID {
		t.Error("rabbitmq dequeue task name != queued task name")
	}

	if string(dequeuedTask.Payload) != string(task.Payload) {
		t.Error("rabbitmq dequeue task payload != queued task payload")
	}

	queue.Close()
}
