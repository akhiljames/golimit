package golimit_test

import (
	"fmt"
	"sync"
	"time"

	"github.com/akhiljames/golimit"
	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

func (r *Redis) Eval(script string, keys []string, args ...interface{}) (interface{}, error) {
	return r.client.Eval(script, keys, args...).Result()
}

func ExampleTake() {
	limiter := golimit.New(
		&Redis{redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})},
		"golimit:test1",
		&golimit.Bucket{
			Interval: 1 * time.Second,
			Quantum:  5,
			Capacity: 10,
		},
	)
	if ok, _ := limiter.Take(1); ok {
		fmt.Println("PASS")
	} else {
		fmt.Println("DROP")
	}
	// Output:
	// PASS
}

func ExampleTake_concurrency() {
	concurrency := 5
	limiter := golimit.New(
		&Redis{redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})},
		"golimit:test2",
		&golimit.Bucket{
			Interval: 1 * time.Second,
			Quantum:  5,
			Capacity: 10,
		},
	)

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			if ok, _ := limiter.Take(1); ok {
				fmt.Println("PASS")
			} else {
				fmt.Println("DROP")
			}
			wg.Done()
		}()
	}
	wg.Wait()
	// Output:
	// PASS
	// PASS
	// PASS
	// PASS
	// PASS
}
