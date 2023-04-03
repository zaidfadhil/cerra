package goatq

import "time"

type Task struct {
	Name    string
	Payload []byte

	retry   int
	timeout time.Duration
	delay   time.Duration
}

func NewTask(name string, payload []byte) *Task {
	return &Task{
		Name:    name,
		Payload: payload,
	}
}

// func (t *Task) Delay(time time.Duration) *Task {
// 	t.delay = time
// 	return t
// }
//
// func (t *Task) Timeout(time time.Duration) *Task {
// 	t.timeout = time
// 	return t
// }
//
// func (t *Task) Retry(num int) *Task {
// 	t.retry = num
// 	return t
// }
