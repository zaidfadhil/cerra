package goatq

import (
	"context"
	"log"
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

var startOnce = sync.Once{}

func NewManager() *Manager {
	var manager Manager
	//var startOnce = sync.Once{}
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

	pid := syscall.Getpid()

	for {
		select {
		case sig := <-c:
			switch sig {
			case syscall.SIGINT, syscall.SIGTSTP, syscall.SIGTERM:
				log.Printf("pid %v shutdown", pid)
				m.shutdown()
				return
			default:
				log.Printf("pid %v sig %v", pid, sig)
			}
		case <-ctx.Done():
			m.shutdown()
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
					log.Printf("shutdown func error %v", err)
				}
			})
		}(f)
	}

	go func() {
		m.runningWaitGroup.Wait()
		m.mutex.Lock()
		m.doneCtxCancel()
		m.mutex.Unlock()
	}()
}

func (m *Manager) OnShutdown(f func() error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.runAtShutdown = append(m.runAtShutdown, f)
}
