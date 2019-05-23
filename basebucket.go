package golimit

import (
	"sync"
	"time"
)

// Config is the bucket configuration.
type Config struct {
	Interval time.Duration // the interval between each addition of one token

	// the capacity of the bucket
	Capacity int64
}

// baseBucket basic structure.
type baseBucket struct {
	mu     sync.RWMutex
	config *Config
}

// Config returns the bucket configuration in a concurrency-safe way.
func (b *baseBucket) Config() Config {
	b.mu.RLock()
	config := *b.config
	b.mu.RUnlock()
	return config
}

// SetConfig updates the bucket configuration in a concurrency-safe way.
func (b *baseBucket) SetConfig(config *Config) {
	b.mu.Lock()
	b.config = config
	b.mu.Unlock()
}
