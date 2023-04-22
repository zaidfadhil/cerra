package rabbitmq_test

import (
	"testing"

	"github.com/zaidfadhil/goatq"
	"github.com/zaidfadhil/goatq/rabbitmq"
)

func TestRabbitmqEnqueue(t *testing.T) {
	backend := rabbitmq.New(rabbitmq.Options{
		Address: "amqp://user:pass@localhost:5672",
	})
	defer backend.Close()

	task := &goatq.Task{
		Name:    "test_task",
		Payload: []byte("test_payload"),
	}

	err := backend.Enqueue(task)
	if err != nil {
		t.Errorf("rabbitmq enqueu error: %v", err)
	}
}

func TestRabbitmqDequeue(t *testing.T) {
	backend := rabbitmq.New(rabbitmq.Options{
		Address: "amqp://user:pass@localhost:5672",
	})
	defer backend.Close()

	task := &goatq.Task{
		Name:    "test_task",
		Payload: []byte("test_payload"),
	}

	err := backend.Enqueue(task)
	if err != nil {
		t.Errorf("rabbitmq enqueu error: %v", err)
	}

	dequeuedTask, err := backend.Dequeue()
	if err != nil {
		t.Errorf("rabbitmq dequeu error: %v", err)
	}

	if dequeuedTask.Name != task.Name {
		t.Error("rabbitmq dequeue task name != queued task name")
	}

	if string(dequeuedTask.Payload) != string(task.Payload) {
		t.Error("rabbitmq dequeue task payload != queued task payload")
	}
}

func TestRabbitmqClose(t *testing.T) {
	backend := rabbitmq.New(rabbitmq.Options{
		Address: "amqp://user:pass@localhost:5672",
	})
	err := backend.Close()
	if err != nil {
		t.Errorf("rabbitmq close connection error: %v", err)
	}
}
