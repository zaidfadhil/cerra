package goatq

import (
	"sync"
)

type inMemoryBackend struct {
	sync.Mutex
	tasks []*Task
}

func NewInMemoryBackend() *inMemoryBackend {
	return &inMemoryBackend{}
}

func (b *inMemoryBackend) Push(task *Task) error {
	b.Lock()
	defer b.Unlock()

	b.tasks = append(b.tasks, task)
	return nil
}

func (b *inMemoryBackend) Request() (*Task, error) {
	if len(b.tasks) == 0 {
		return nil, ErrEmtpyQueue
	}

	b.Lock()
	defer b.Unlock()

	data := b.tasks[:1][0]
	b.tasks = b.tasks[1:]

	return data, nil
}

func (b *inMemoryBackend) Close() error {
	return nil
}
