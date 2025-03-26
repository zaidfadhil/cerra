# Cerra
[![Test](https://github.com/zaidfadhil/cerra/actions/workflows/test.yaml/badge.svg?branch=main)](https://github.com/zaidfadhil/cerra/actions/workflows/test.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/zaidfadhil/cerra.svg)](https://pkg.go.dev/github.com/zaidfadhil/cerra)
[![Go Report Card](https://goreportcard.com/badge/github.com/zaidfadhil/cerra)](https://goreportcard.com/report/github.com/zaidfadhil/cerra)

Cerra is a simple task queue library in Go that supports in-memory, Redis, and rabbitmq (amqp) backends.

## Features

* [x] Support In-Memory
* [x] Support Redis
* [x] Support rabbitmq (amqp)

Resources:

- [Examples](https://github.com/zaidfadhil/cerra/tree/main/_example)
- [Reference](https://pkg.go.dev/github.com/zaidfadhil/cerra)

## Requirements

Make sure you have Go installed. Version 1.18 or higher.

## Installation

To install cerra, use go get:
```bash
go get github.com/zaidfadhil/cerra
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/zaidfadhil/cerra"
)

func main() {
	// Create a new queue with the in-memory backend
	// Set Max number of workers. 0 for max number of workers (runtime.NumCPU() * 2)
        queue := cerra.NewQueue(cerra.NewInMemoryBackend(), 0)

	// Update max number of workers.
	queue.UpdateMaxWorkerNum(5)

	// Add a handler function
	queue.AddHandler(func(ctx context.Context, task *cerra.Task) error {
		fmt.Printf("Received task with ID %s and payload %v\n", task.ID, task.Payload)
		return nil
	})

	// Start the queue
	queue.Start()

	// Enqueue some tasks
	for i := 0; i < 10; i++ {
		task := cerra.NewTask([]byte(fmt.Sprint(i)))
		err := queue.Enqueue(task) 
                if err != nil {
			fmt.Printf("error enqueueing task: %v\n", err)
		}
		fmt.Println("Enqueue", i)
	}

	// Wait for the tasks to be processed
	time.Sleep(3 * time.Second)

	// Close the queue
	queue.Close()
}
```

### More Backends

to use Redis as a backend for the queues, just replace the in-memory backend with redis

```go
// Create Redis Backend
backend := redis.New(redis.Options{
	Address:  "localhost:6379",
	Password: "redis",
	Stream:   "cerra",
	Group:    "cerra",
	Consumer: "cerra",
})

// Create a new queue
queue := cerra.NewQueue(backend, 0)
```

and the same for using rabbitmq with amqp

```go
// Create amqp Backend
backend := amqp.New(amqp.Options{
	Address:      "amqp://user:pass@localhost:5672/",
	Queue:        "cerra",
	ExchangeName: "cerra-exchange",
	ExchangeType: "direct",
	RoutingKey:   "cerra-key",
})

// Create a new queue
queue := cerra.NewQueue(backend, 0)
```

## License

Cerra is licensed under the [MIT License](https://github.com/zaidfadhil/cerra/blob/master/LICENSE).
