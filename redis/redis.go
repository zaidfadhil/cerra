package redisq

type RedisBackend struct{}

// func NewRedisQueue(client *redis.Client) *Queue {
// 	return &Queue{
// 		Backend: NewRedisBackend(client),
// 	}
// }
//
// func NewRedisBackend(client *redis.Client) *RedisBackend {
// 	return &RedisBackend{client: client}
// }
//
// func (b *RedisBackend) Push(task *Task) error {
// 	return b.client.RPush(context.Background(), task.Name, task.Payload).Err()
// }
//
// func (b *RedisBackend) Run(handler func(context.Context, *Task) error) error {
// 	for {
// 		// task, err := b.client.BLPop(context.Background(), 0, task.Name).Result()
// 		// if err != nil {
// 		// 	return err
// 		// }
// 		//
// 		// if err := handler(context.Background(), NewTask(task.Name, []byte(task.Payload))); err != nil {
// 		// 	return err
// 		// }
// 	}
// }
//
// func (b *RedisBackend) Close() error {
// 	return b.client.Close()
// }
