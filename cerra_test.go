package cerra_test

import (
	"context"
	"errors"
	"sync/atomic"
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

	queue.AddHandler(func(ctx context.Context, t *cerra.Task) error {
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
}

func TestTaskRetry(t *testing.T) {
	queue := cerra.NewQueue(cerra.NewInMemoryBackend(), 1)

	var taskRetryCount int32

	queue.AddHandler(func(ctx context.Context, t *cerra.Task) error {
		atomic.AddInt32(&taskRetryCount, 1)
		return errors.New("retry error")
	})

	queue.Start()
	defer queue.Close()

	task := &cerra.Task{
		ID:         "test_task",
		Payload:    []byte("test_payload"),
		RetryLimit: 5,
	}

	err := queue.Enqueue(task)
	if err != nil {
		t.Errorf("enqueue error: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	if atomic.LoadInt32(&taskRetryCount) != 6 {
		t.Errorf("wrong task retry count. %v", taskRetryCount)
	}
}
