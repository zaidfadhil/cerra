package cerra

import (
	"encoding/json"
	"time"
)

type Task struct {
	ID         string        `json:"id"`
	Payload    []byte        `json:"payload"`
	Timeout    time.Duration `json:"timeout"`
	RetryLimit int           `json:"retry_limit"`
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
		Timeout: 60 * time.Second,
	}
}

func (t *Task) SetID(id string) {
	t.ID = id
}

func (t *Task) SetRetryLimit(limit int) {
	if limit >= 0 {
		t.RetryLimit = limit
	}
}

func (t *Task) SetTimeout(timeout time.Duration) {
	t.Timeout = timeout
}

func (t *Task) Encode() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Task) ToMap() map[string]any {
	return map[string]any{
		"id":          t.ID,
		"payload":     t.Payload,
		"timeout":     t.Timeout,
		"retry_limit": t.RetryLimit,
		"retry_count": t.RetryCount,
	}
}
