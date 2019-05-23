package golimit_test

import (
	"fmt"
	"strings"
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

// EvalSha Optimize Lua script execution
func (r *Redis) EvalSha(sha1 string, keys []string, args ...interface{}) (interface{}, error, bool) {
	result, err := r.client.EvalSha(sha1, keys, args...).Result()
	noScript := err != nil && strings.HasPrefix(err.Error(), "NOSCRIPT ")
	return result, err, noScript
}

func ExampleRateLimiter_Take() {
	tb := golimit.NewRateLimiter(
		&Redis{redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})},
		"golimit:example",
		&golimit.Config{
			Interval: 1 * time.Second / 2,
			Capacity: 5,
		},
	)
	if ok, err := tb.Take(1); ok {
		fmt.Println("PASS")
	} else {
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("DROP")
	}
	// Output:
	// PASS
}

func ExampleTake_concurrency() {
	concurrency := 5
	limiter := golimit.NewRateLimiter(
		&Redis{redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})},
		"golimit:concurrency",
		&golimit.Config{
			Interval: 1 * time.Second,
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
