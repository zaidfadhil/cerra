package goatq_test

import (
	"testing"

	"github.com/zaidfadhil/goatq"
)

func TestInMemEnqueue(t *testing.T) {
	backend := goatq.NewInMemoryBackend()
	task := &goatq.Task{
		Name:    "test_task",
		Payload: []byte("test_payload"),
	}

	err := backend.Enqueue(task)
	if err != nil {
		t.Errorf("redisq enqueu error: %v", err)
	}
}

func TestInMemDequeue(t *testing.T) {
	backend := goatq.NewInMemoryBackend()

	task := &goatq.Task{
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
