# Cerra
Cerra is a simple task queue library in Go that supports in-memory, Redis, and RabbitMQ backends.

## Features

* [x] Support In-Memory with buffer channels
* [x] Support Redis
* [x] Support RabbitMQ

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
	// Set Max number of workers. default number is (runtime.NumCPU() * 2)
	queue := cerra.NewQueue(cerra.NewInMemoryBackend(), 2)

	// Update max number of workers.
	queue.UpdateMaxWorkerNum(5)

	// Add a handler function to the queue
	queue.AddHandler(func(ctx context.Context, task *cerra.Task) error {
		fmt.Printf("Received task with Name %s and payload %v\n", task.Name, task.Payload)
		return nil
	})

	// Start the queue
	queue.Start()

	// Enqueue some tasks
	for i := 0; i < 10; i++ {
		task := cerra.NewTask(fmt.Sprint(i), []byte(fmt.Sprint(i)))
		if err := queue.Enqueue(task); err != nil {
			fmt.Printf("Error enqueueing task: %v\n", err)
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

to use Redis as a backend for the queues, just replace the in-memory backend with redisq

```go
// Create Redis Backend
backend := redisq.New(redisq.Options{
	Address:  "localhost:6379",
	Password: "redis",
	Stream:   "cerra",
	Group:    "cerra",
	Consumer: "cerra",
})

// Create a new queue
queue := cerra.NewQueue(backend, 0)
```

and the same for using RabbitMQ

```go
// Create RabbitMQ Backend
backend := rabbitmq.New(rabbitmq.Options{
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
