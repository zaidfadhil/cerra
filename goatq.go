package goatq

import (
	"context"
	"errors"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Backend interface {
	Enqueue(task *Task) error
	Dequeue() (*Task, error)
	Close() error
}

type Queue struct {
	sync.Mutex
	Backend      Backend
	group        *routineGroup
	quit         chan struct{}
	ready        chan struct{}
	stop         sync.Once
	maxWorkerNum int
	handleFuncs  []func(context.Context, *Task) error

	activeWorkers uint32
}

func NewQueue(backend Backend) *Queue {
	return &Queue{
		Backend:      backend,
		group:        newRoutineGroup(),
		quit:         make(chan struct{}),
		ready:        make(chan struct{}, 1),
		maxWorkerNum: runtime.NumCPU(),
	}
}

func (q *Queue) Enqueue(t *Task) error {
	return q.Backend.Enqueue(t)
}

func (q *Queue) UpdateMaxWorkerNum(num int) {
	if num != 0 {
		q.maxWorkerNum = num
		q.schedule()
	}
}

func (q *Queue) Close() {
	q.stop.Do(func() {
		err := q.Backend.Close()
		if err != nil {
			log.Println(err)
		}
		close(q.quit)
	})
	q.group.Wait()
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
				task, err := q.Backend.Dequeue()
				if err != nil {
					select {
					case <-q.quit:
						if !errors.Is(err, ErrEmtpyQueue) {
							close(tasks)
							return
						}
					case <-time.After(time.Second):
					}
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

		atomic.AddUint32(&q.activeWorkers, 1)
		q.group.Run(func() {
			q.runFunc(ctx, task)
		})
	}
}

func (q *Queue) schedule() {
	q.Lock()
	defer q.Unlock()

	if int(q.activeWorkers) >= q.maxWorkerNum {
		return
	}

	select {
	case q.ready <- struct{}{}:
	default:
	}
}

func (q *Queue) runFunc(ctx context.Context, t *Task) {
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer func() {
		atomic.AddUint32(&q.activeWorkers, ^uint32(0))
		q.schedule()
		cancel()
	}()

	for _, f := range q.handleFuncs {
		if err := f(ctx, t); err != nil {
			log.Printf("internal error: %v", err)
		}
	}
}
