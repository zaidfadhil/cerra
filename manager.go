package goatq

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Manager struct {
	mutex             *sync.RWMutex
	shutdownCtx       context.Context
	shutdownCtxCancel context.CancelFunc
	doneCtx           context.Context
	doneCtxCancel     context.CancelFunc
	runningWaitGroup  *routineGroup
	runAtShutdown     []func() error
}

func NewManager() *Manager {
	var manager Manager
	var startOnce = sync.Once{}
	startOnce.Do(func() {
		manager = Manager{
			mutex:            &sync.RWMutex{},
			runningWaitGroup: newRoutineGroup(),
		}

		ctx := context.Background()

		manager.shutdownCtx, manager.shutdownCtxCancel = context.WithCancel(ctx)
		manager.doneCtx, manager.doneCtxCancel = context.WithCancel(context.Background())

		go manager.waitSignal(ctx)
	})
	return &manager
}

func (m *Manager) waitSignal(ctx context.Context) {
	c := make(chan os.Signal, 1)
	signal.Notify(
		c,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
	)
	defer signal.Stop(c)
	for {
		select {
		case sig := <-c:
			switch sig {
			case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
				m.shutdown()
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (m *Manager) Done() <-chan struct{} {
	return m.doneCtx.Done()
}

func (m *Manager) shutdown() {
	m.shutdownCtxCancel()

	for _, f := range m.runAtShutdown {
		func(run func() error) {
			m.runningWaitGroup.Run(func() {
				if err := f(); err != nil {
					fmt.Println(err)
				}
			})
		}(f)
	}

	go func() {
		m.mutex.Lock()
		m.doneCtxCancel()
		m.mutex.Unlock()
	}()
}

func (m *Manager) OnShutdown(f func() error) {
	m.mutex.Lock()
	m.runAtShutdown = append(m.runAtShutdown, f)
	m.mutex.Unlock()
}
