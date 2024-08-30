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
	retryCount   int

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
		retryCount:   0,
	}
}

func (q *Queue) Enqueue(t *Task) error {
	if t.RetryCount <= 0 && q.retryCount > 0 {
		t.SetRetry(q.retryCount)
	}

	return q.Backend.Enqueue(t)
}

func (q *Queue) UpdateMaxWorkerNum(num int) {
	if num != 0 {
		q.maxWorkerNum = num
		q.schedule()
	}
}

func (q *Queue) SetRetryCount(num int) {
	if num >= 0 {
		q.retryCount = num
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
	ctx, cancel := context.WithTimeout(ctx, t.Timeout)
	defer func() {
		atomic.AddUint32(&q.activeWorkers, ^uint32(0))
		q.schedule()
		cancel()
	}()

	for _, f := range q.handleFuncs {
		if t.RetryCount <= 0 && q.retryCount > 0 {
			t.SetRetry(q.retryCount)
		}

		err := f(ctx, t)
		if err != nil {
			log.Printf("internal error: %v", err)
			if t.RetryCount > 0 {
				for i := 0; i < t.RetryCount; i++ {
					if err := f(ctx, t); err != nil {
						log.Printf("internal error: %v, retry_count: %v", err, i+1)
					}
				}
			}
		}
	}
}
