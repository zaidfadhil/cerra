package redisq_test

import (
	"testing"

	"github.com/zaidfadhil/cerra"
	"github.com/zaidfadhil/cerra/redisq"
)

func TestRedisqEnqueue(t *testing.T) {
	backend := redisq.New(redisq.Options{
		Address: "localhost:6379",
	})
	defer backend.Close()

	task := &cerra.Task{
		Name:    "test_task",
		Payload: []byte("test_payload"),
	}

	err := backend.Enqueue(task)
	if err != nil {
		t.Errorf("redisq enqueu error: %v", err)
	}
}

func TestRedisqDequeue(t *testing.T) {
	backend := redisq.New(redisq.Options{
		Address: "localhost:6379",
	})
	defer backend.Close()

	task := &cerra.Task{
		Name:    "test_task",
		Payload: []byte("test_payload"),
	}

	err := backend.Enqueue(task)
	if err != nil {
		t.Errorf("redisq enqueu error: %v", err)
	}

	dequeuedTask, err := backend.Dequeue()
	if err != nil {
		t.Errorf("redisq dequeu error: %v", err)
	}

	if dequeuedTask.Name != task.Name {
		t.Error("redisq dequeue task name != queued task name")
	}

	if string(dequeuedTask.Payload) != string(task.Payload) {
		t.Error("redisq dequeue task payload != queued task payload")
	}
}

func TestRedisqClose(t *testing.T) {
	backend := redisq.New(redisq.Options{
		Address: "localhost:6379",
	})
	err := backend.Close()
	if err != nil {
		t.Errorf("redisq close connection error: %v", err)
	}
}
