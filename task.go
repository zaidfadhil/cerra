package cerra

import (
	"encoding/json"
	"time"
)

type Task struct {
	ID      string        `json:"id"`
	Payload []byte        `json:"payload"`
	Timeout time.Duration `json:"timeout"`
}

func NewTask(payload []byte) *Task {
	return &Task{
		Payload: payload,
		Timeout: 60 * time.Minute,
	}
}

func NewTaskWithID(id string, payload []byte) *Task {
	return &Task{
		ID:      id,
		Payload: payload,
		Timeout: 60 * time.Minute,
	}
}

func (t *Task) Encode() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Task) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":      t.ID,
		"payload": t.Payload,
		"timeout": t.Timeout,
	}
}

func (t *Task) AddTimeout(time time.Duration) *Task {
	t.Timeout = time
	return t
}
