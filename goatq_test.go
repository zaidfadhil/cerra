package goatq_test

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/zaidfadhil/goatq"
)

type Model struct {
	Name string
	Age  int
}

func newTask(name string, age int) (*goatq.Task, error) {
	payload, err := json.Marshal(&Model{Name: name, Age: age})
	if err != nil {
		return nil, err
	}
	fmt.Println("set:", payload)
	return goatq.NewTask("yo:queue", payload), nil
}

func handleTask(ctx context.Context, t *goatq.Task) error {
	var model Model
	if err := json.Unmarshal(t.Payload, &model); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}
	fmt.Println("get:", model)
	return nil
}

func Test(t *testing.T) {

	PrintMemUsage()

	queue := goatq.NewQueue(goatq.NewInMemoryBackend())
	defer queue.Close()

	for i := 0; i < 100; i++ {
		task, err := newTask("test-queue", i)
		if err != nil {
			t.Error(err)
		}

		err = queue.Push(task)
		if err != nil {
			t.Error(err)
		}
	}

	PrintMemUsage()

	queue.Handle(handleTask)

	time.Sleep(1 * time.Second)

	PrintMemUsage()
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
