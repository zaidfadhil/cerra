package goatq

import (
	"sync"
)

type inMemoryBackend struct {
	sync.RWMutex
	tasks []*Task

	head  int
	tail  int
	count int
}

func NewInMemoryBackend() *inMemoryBackend {
	return &inMemoryBackend{
		tasks: make([]*Task, 1),
	}
}

func (b *inMemoryBackend) Push(task *Task) error {
	b.Lock()
	defer b.Unlock()

	if b.count == len(b.tasks) {
		b.resize(b.count * 2)
	}

	b.tasks[b.tail] = task
	b.tail = (b.tail + 1) % len(b.tasks)
	b.count++

	return nil
}

func (b *inMemoryBackend) Request() (*Task, error) {
	if len(b.tasks) == 0 {
		return nil, ErrEmtpyQueue
	}

	b.Lock()
	defer b.Unlock()

	data := b.tasks[b.head]
	b.tasks[b.head] = nil
	b.head = (b.head + 1) % len(b.tasks)
	b.count--

	if n := len(b.tasks) / 2; n > 1 && b.count <= n {
		b.resize(n)
	}

	return data, nil
}

func (b *inMemoryBackend) Close() error {
	return nil
}

func (b *inMemoryBackend) Size() int {
	b.Lock()
	defer b.Unlock()
	return len(b.tasks)
}

func (b *inMemoryBackend) resize(size int) {
	nodes := make([]*Task, size)

	if b.head < b.tail {
		copy(nodes, b.tasks[b.head:b.tail])
	} else {
		copy(nodes, b.tasks[b.head:])
		copy(nodes[len(b.tasks)-b.head:], b.tasks[:b.tail])
	}

	b.tail = b.count % size
	b.head = 0
	b.tasks = nodes
}
