package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	redis "github.com/go-redis/redis/v9"
)

var client *redis.Client = redis.NewClient(&redis.Options{
	Network: "tcp",
	Addr:    "localhost:6379",
})

var ctx context.Context = context.Background()

const (
	release_script = `
		if redis.call('get', KEYS[1]) == ARGV[1] then
			return redis.call('del', KEYS[1])
		else
			return 0
		end
	`
	resource_name = "queue"
	ttl           = 10
	num_consumers = 3
)

type Message struct {
	id        int
	message   string
	processed bool
}

type Queue struct {
	events       []*Message
	messageCount int
	owner        string
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

	for {
		res := client.SetNX(ctx, resource_name, client_id, time.Duration(time.Duration.Seconds(10)))
		if res.Val() == true {
			fmt.Printf("Acquiring Lock :: %v\n", client_id)
			return 1
		} else {
			continue
		}
	}

}

func ReleaseLock(resource_name, client_id string) bool {
	release_lock := redis.NewScript(release_script)
	key := []string{resource_name}
	fmt.Printf("Releasing Lock :: %v\n", client_id)
	_, err := release_lock.Run(ctx, client, key, client_id).Int()
	if err != nil {
		fmt.Printf("Can not release lock: %v\n", err)
	}
	return true
}

func (q *Queue) _ProcessMessage(client_id string) {
	for {
		if q.messageCount == 0 {
			break
		}
		num := AcquireLock(resource_name, client_id, ttl)
		if &num != nil && q.messageCount > 0 {
			// fmt.Printf("Client : %v for message : %v\n", client_id, q.events[q.messageCount-1].id)
			if q.messageCount > 0 && q.events[q.messageCount-1].processed != true {
				fmt.Printf("Queue message processed by: %s => id: %v, message: %v\n",
					client_id, q.events[q.messageCount-1].id, q.events[q.messageCount-1].message)

				q.events[q.messageCount-1].processed = true
				q.messageCount--
			}
		}
		ReleaseLock(resource_name, client_id)
	}

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
			q._ProcessMessage(fmt.Sprintf("client_%v", i%num_consumers))
		}()
	}

	wg.Wait()

}
