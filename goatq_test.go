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
	Title string
	Num   int
}

func newTask(title string, num int) (*goatq.Task, error) {
	payload, err := json.Marshal(&Model{Title: title, Num: num})
	if err != nil {
		return nil, err
	}
	return goatq.NewTask("yo:queue", payload), nil
}

func handleTask(ctx context.Context, t *goatq.Task) error {
	var model Model
	if err := json.Unmarshal(t.Payload, &model); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}
	return nil
}

func Test(t *testing.T) {
	PrintMemUsage()

	queue := goatq.NewQueue(goatq.NewInMemoryBackend())

	for i := 0; i < 1000000; i++ {
		task, err := newTask("test-queue", i)
		if err != nil {
			fmt.Println(err)
		}

		err = queue.Enqueue(task)
		if err != nil {
			fmt.Println(err)
		}
	}

	queue.AddHandler(handleTask)

	queue.Start()

	PrintMemUsage()

	time.Sleep(2 * time.Second)

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
