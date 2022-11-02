package main

import (
	"context"
	"fmt"
	"sync"

	redis "github.com/go-redis/redis/v9"
)

var client *redis.Client = redis.NewClient(&redis.Options{
	Network: "tcp",
	Addr:    "localhost:6379",
})

var ctx context.Context = context.Background()

const (
	acquire_script = `
		return redis.call('SET', KEYS[1], ARGV[1], 'NX', 'PX', ARGV[2])
	`

	release_script = `
		if redis.call('GET', KEYS[1]) == ARGV[1] then
			return redis.call('DEL', KEYS[1])
		else
			return 0
		end
	`
	resource_name = "queue"
	ttl           = 10
	num_consumers = 3
)

type Message struct {
	id      int
	message string
}

type Queue struct {
	events       []*Message
	messageCount int
}

var q *Queue

func InitQueue() {
	events := []*Message{
		&Message{id: 1, message: "message 1"},
		&Message{id: 2, message: "message 2"},
		&Message{id: 3, message: "message 3"},
		&Message{id: 4, message: "message 4"},
	}

	q = &Queue{
		events:       events,
		messageCount: 4,
	}
}

func AcquireLock(resource_name, client_id string, ttl int) int {

	acquire_lock := redis.NewScript(acquire_script)
	key := []string{resource_name}
	for {
		num, _ := acquire_lock.Run(ctx, client, key, client_id, ttl).Int()
		if &num != nil {
			fmt.Printf("Lock acquired by: %v\n", client_id)
			return num
		} else {
			continue
		}
	}
}

func ReleaseLock(resource_name, client_id string) int {
	release_lock := redis.NewScript(release_script)
	key := []string{resource_name}

	num, err := release_lock.Run(ctx, client, key, client_id).Int()
	if err != nil {
		fmt.Printf("Can not release lock: %v\n", err)
	}
	fmt.Printf("Lock released: %v\n", client_id)
	return num
}

func (q *Queue) ProcessMessage(client_id string) {
	for q.messageCount > 0 {
		num := AcquireLock(resource_name, client_id, ttl)
		if &num != nil && q.messageCount > 0 {
			fmt.Printf("Queue message processed by: %s => id: %v, message: %v\n", client_id, q.events[q.messageCount-1].id, q.events[q.messageCount-1].message)
			q.messageCount--
		}
	}

	ReleaseLock(resource_name, client_id)
	fmt.Printf("All message processed : %v\n", client_id)

}

func main() {
	var wg sync.WaitGroup

	InitQueue()

	var i int

	for i = 0; i < num_consumers; i++ {
		wg.Add(1)

		i := i

		go func() {
			defer wg.Done()
			q.ProcessMessage(fmt.Sprintf("client_%v", i%num_consumers))
		}()
	}

	wg.Wait()

}
