package goatq

import (
	"context"
	"time"
)

type Backend interface {
	Push(task *Task) error
	Request() (*Task, error)
	Close() error
}

type Queue struct {
	Backend Backend

	routineGroup *routineGroup
}

func NewQueue(client Backend) *Queue {
	return &Queue{
		Backend:      client,
		routineGroup: newRoutineGroup(),
	}
}

func (q *Queue) Close() error {
	return q.Backend.Close()
}

func (q *Queue) Push(t *Task) error {
	return q.Backend.Push(t)
}

func (q *Queue) Handle(handler func(context.Context, *Task) error) {
	if handler == nil {
		panic("goatq: nil handler")
	}

	ctx := context.Background()

	q.routineGroup.Run(func() {
	loop:
		for {
			task, err := q.Backend.Request()
			if err != nil {
				return
			}

			handler(ctx, task)

			select {
			case <-time.After(task.delay):
			case <-ctx.Done(): // timeout reached
				//err = ctx.Err()
				break loop
			}
		}
	})
}
