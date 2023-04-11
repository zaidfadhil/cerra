package goatq

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Backend interface {
	Enqueue(task *Task) error
	Request() (*Task, error)
	Close() error
}

type Queue struct {
	sync.Mutex
	Backend Backend

	group *routineGroup

	quit  chan struct{}
	ready chan struct{}
	stop  sync.Once

	handleFuncs []func(context.Context, *Task) error
}

func NewQueue(client Backend) *Queue {
	return &Queue{
		Backend: client,
		group:   newRoutineGroup(),
		quit:    make(chan struct{}),
		ready:   make(chan struct{}, 1),
	}
}

func (q *Queue) Enqueue(t *Task) error {
	return q.Backend.Enqueue(t)
}

func (q *Queue) Close() error {
	q.stop.Do(func() {
		q.Backend.Close()
		close(q.quit)
		q.group.Wait()
	})
	return nil
}

func (q *Queue) AddHandler(handler func(context.Context, *Task) error) {
	q.Lock()
	q.handleFuncs = append(q.handleFuncs, handler)
	q.Unlock()
}

func (q *Queue) Start() {
	q.group.Run(func() {
		q.start()
	})
}

func (q *Queue) start() {
	tasks := make(chan *Task, 1)

	ctx := context.Background()

	for {
		q.schedule()

		select {
		case <-q.ready:
		case <-q.quit:
			return
		}

		q.group.Run(func() {
			for {
				task, err := q.Backend.Request()
				if err != nil {
					return
				}

				if task != nil {
					tasks <- task
					return
				}

				select {
				case <-q.quit:
					if !errors.Is(err, ErrEmtpyQueue) {
						close(tasks)
						return
					}
				default:
				}
			}
		})

		task, ok := <-tasks
		if !ok {
			return
		}

		q.group.Run(func() {
			q.runFunc(ctx, task)
		})
	}
}

func (q *Queue) schedule() {
	q.Lock()
	defer q.Unlock()

	select {
	case q.ready <- struct{}{}:
	default:
	}
}

func (q *Queue) runFunc(ctx context.Context, t *Task) {
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer func() {
		cancel()
	}()

	q.group.Run(func() {
		for _, f := range q.handleFuncs {
			if err := f(ctx, t); err != nil {
				fmt.Println("internal error", err)
			}
		}
	})
}
