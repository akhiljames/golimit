# golimit

A simple redis backed rate limiter for distributed systems based on token bucket algorithm.

## Install
```shell
go get "github.com/akhiljames/golimit"
```

## Usage
```go
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/akhiljames/golimit"
	"github.com/go-redis/redis"
)

// Redis client. You can use your own redis client
type Redis struct {
	client *redis.Client
}

// Eval exec script
func (r *Redis) Eval(script string, keys []string, args ...interface{}) (interface{}, error) {
	return r.client.Eval(script, keys, args...).Result()
}

// EvalSha Optimize Lua script execution
func (r *Redis) EvalSha(sha1 string, keys []string, args ...interface{}) (interface{}, error, bool) {
	result, err := r.client.EvalSha(sha1, keys, args...).Result()
	noScript := err != nil && strings.HasPrefix(err.Error(), "NOSCRIPT ")
	return result, err, noScript
}

func main() {
	// create a new rate limiter
	tb := golimit.NewRateLimiter(
		&Redis{redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})},
		"golimit:example", // this is a unique key for redis
		&golimit.Config{
			Interval: 1 * time.Second, // interval for rate limit
			Capacity: 5,               // capacity to be processed in the interval
			// this will rate limt as per - capacity amount in interval time
		},
	)
	if ok, err := tb.Take(1); ok { // usage of operation can be taken. for eg. heavy db ops
		// serve the user request
	} else {
		if err != nil {
			fmt.Println(err.Error())
		}
		// reject the user request
	}
}
```
