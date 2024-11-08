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
	if err := queue.Enqueue(task); err != nil {
		t.Errorf("redis enqueue error: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
}

func TestRedisDequeue(t *testing.T) {
	backend := redisq.New(redisq.Options{
		Address: "localhost:6379",
	})
	queue := cerra.NewQueue(backend, 1)
	queue.Start()
	defer queue.Close()

	time.Sleep(50 * time.Millisecond)
	task := &cerra.Task{
		ID:      "test_task",
		Payload: []byte("test_payload"),
	}
	err := queue.Enqueue(task)
	if err != nil {
		t.Errorf("redis enqueue error: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	queue.AddHandler(func(ctx context.Context, tt *cerra.Task) error {
		if tt == nil {
			t.Error("handler was not called")
		}
		return nil
	})

	time.Sleep(200 * time.Millisecond)
}
