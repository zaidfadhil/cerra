package cerra_test

import (
	"testing"

	"github.com/zaidfadhil/cerra"
)

func TestInMemEnqueue(t *testing.T) {
	backend := cerra.NewInMemoryBackend()
	task := &cerra.Task{
		Payload: []byte("test_payload"),
	}

	err := backend.Enqueue(task)
	if err != nil {
		t.Errorf("redisq enqueu error: %v", err)
	}
}

func TestInMemDequeue(t *testing.T) {
	backend := cerra.NewInMemoryBackend()

	task := &cerra.Task{
		ID:      "id",
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

	if dequeuedTask.ID != task.ID {
		t.Error("redisq dequeue task name != queued task name")
	}

	if string(dequeuedTask.Payload) != string(task.Payload) {
		t.Error("redisq dequeue task payload != queued task payload")
	}
}
