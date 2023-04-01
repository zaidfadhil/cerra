package goatq_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

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

	client := goatq.NewQueue()
	defer client.Close()

	task, err := newTask("test-queue", 20)
	if err != nil {
		t.Error(err)
	}

	err = client.Push(task)
	if err != nil {
		t.Error(err)
	}

	client.Handle("test-queue", handleTask)
	err = client.Run()
	if err != nil {
		t.Error(err)
	}

}
