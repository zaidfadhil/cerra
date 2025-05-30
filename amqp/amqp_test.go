package amqp_test

import (
	"context"
	"testing"
	"time"

	"github.com/zaidfadhil/cerra"
	"github.com/zaidfadhil/cerra/amqp"
)

func TestAmqpEnqueue(t *testing.T) {
	backend := amqp.New(amqp.Options{
		Address: "amqp://user:pass@localhost:5672",
	})
	queue := cerra.NewQueue(backend, 1)
	defer queue.Close()

	task := &cerra.Task{
		ID:      "test_task",
		Payload: []byte("test_payload"),
	}

	err := queue.Enqueue(task)
	if err != nil {
		t.Errorf("amqp enqueue error: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
}

func TestAmqpDequeue(t *testing.T) {
	backend := amqp.New(amqp.Options{
		Address: "amqp://user:pass@localhost:5672",
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
		t.Errorf("amqp enqueue error: %v", err)
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
