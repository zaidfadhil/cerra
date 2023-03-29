package goatq

// opts := &goatq.RedisOpts{}
// opts := &goatq.InMemOpts{}
// opts := &goatq.RabbitMQOpts{}

// newQueue, err := goatq.RegisterQueue("test-queue", opts)
// err := newQueue.Add({title:"test"}).Send()
// err := newQueue.Add({}).Delay()
// newQueue.Close()

// ctx := context.Background()
// task, err := newQueue.Worker(ctx, workerOpts.WorkersCount(5))
// task.Data // {"title":"test"}
