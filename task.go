package cerra

import (
	"encoding/json"
	"time"
)

type Task struct {
	ID         string        `json:"id"`
	Payload    []byte        `json:"payload"`
	Timeout    time.Duration `json:"timeout"`
	RetryCount int           `json:"retry_count"`
}

func NewTask(payload []byte) *Task {
	return &Task{
		Payload: payload,
		Timeout: 60 * time.Second,
	}
}

func NewTaskWithID(id string, payload []byte) *Task {
	return &Task{
		ID:      id,
		Payload: payload,
		Timeout: 60 * time.Minute,
	}
}

func (t *Task) SetID(id string) {
	t.ID = id
}

func (t *Task) SetRetry(count int) {
	if count >= 0 {
		t.RetryCount = count
	}
}

func (t *Task) SetTimeout(timeout time.Duration) {
	t.Timeout = timeout
}

func (t *Task) Encode() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Task) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":         t.ID,
		"payload":    t.Payload,
		"timeout":    t.Timeout,
		"retryCount": t.RetryCount,
	}
}
