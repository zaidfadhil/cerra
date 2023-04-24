package rabbitmq_test

import (
	"context"
	"testing"
	"time"

	"github.com/zaidfadhil/cerra"
	"github.com/zaidfadhil/cerra/rabbitmq"
)

func TestRabbitmqEnqueue(t *testing.T) {
	backend := rabbitmq.New(rabbitmq.Options{
		Address: "amqp://user:pass@localhost:5672",
	})
	queue := cerra.NewQueue(backend)
	defer queue.Close()

	task := &cerra.Task{
		Name:    "test_task",
		Payload: []byte("test_payload"),
	}

	err := queue.Enqueue(task)
	if err != nil {
		t.Errorf("rabbitmq enqueu error: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
}

func TestRabbitmqDequeue(t *testing.T) {
	backend := rabbitmq.New(rabbitmq.Options{
		Address: "amqp://user:pass@localhost:5672",
	})
	queue := cerra.NewQueue(backend)
	queue.Start()

	task := &cerra.Task{
		Name:    "test_task",
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

	if dequeuedTask.Name != task.Name {
		t.Error("rabbitmq dequeue task name != queued task name")
	}

	if string(dequeuedTask.Payload) != string(task.Payload) {
		t.Error("rabbitmq dequeue task payload != queued task payload")
	}

	queue.Close()
}
