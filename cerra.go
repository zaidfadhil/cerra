package cerra

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

func NewQueue(backend Backend, workers int) *Queue {
	if workers <= 0 {
		workers = runtime.NumCPU() * 2
	}
	return &Queue{
		Backend:      backend,
		group:        newRoutineGroup(),
		quit:         make(chan struct{}),
		ready:        make(chan struct{}, 1),
		maxWorkerNum: workers,
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
		if err := q.Backend.Close(); err != nil {
			log.Println(err)
		}
		close(q.quit)
	})
	q.group.Wait()
}

func (q *Queue) AddHandler(handler func(context.Context, *Task) error) {
	q.Lock()
	defer q.Unlock()

	q.handleFuncs = append(q.handleFuncs, handler)
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

	if int(atomic.LoadUint32(&q.activeWorkers)) >= q.maxWorkerNum {
		return
	}

	select {
	case q.ready <- struct{}{}:
	default:
	}
}

func (q *Queue) runFunc(ctx context.Context, t *Task) {
	ctx, cancel := context.WithTimeout(ctx, t.Timeout)
	defer func() {
		atomic.AddUint32(&q.activeWorkers, ^uint32(0))
		q.schedule()
		cancel()
	}()

	for _, f := range q.handleFuncs {
		if err := f(ctx, t); err != nil {
			if t.RetryCount == 0 {
				log.Printf("cerra task error: %v", err)
			} else {
				log.Printf("cerra task error: %v. retry: %v/%v", err, t.RetryCount, t.RetryLimit)
			}

			if t.RetryLimit > 0 && t.RetryCount != t.RetryLimit {
				t.RetryCount++
				if err = q.Enqueue(t); err != nil {
					log.Printf("cerra task retry error: %v. retries: %v/%v", err, t.RetryCount, t.RetryLimit)
				}
			}
		}
	}
}
