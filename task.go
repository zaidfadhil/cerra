package cerra

import (
	"encoding/json"
	"time"
)

type Task struct {
	Name    string `json:"name"`
	Payload []byte `json:"payload"`

	Timeout time.Duration `json:"timeout"`
}

func NewTask(name string, payload []byte) *Task {
	return &Task{
		Name:    name,
		Payload: payload,
		Timeout: 60 * time.Minute,
	}
}

func (t *Task) Encode() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Task) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":    t.Name,
		"payload": t.Payload,
		"timeout": t.Timeout,
	}
}

func (t *Task) AddTimeout(time time.Duration) *Task {
	t.Timeout = time
	return t
}
