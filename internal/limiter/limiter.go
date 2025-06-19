package limiter

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ntayoane/go-rate-limited/internal/storage"
)

type Config struct {
	RequestsPerSecond int64
	BlockDuration     time.Duration
}

type RateLimiter struct {
	storage storage.Storage
	config  Config
}

func NewRateLimiter(storage storage.Storage, config Config) *RateLimiter {
	return &RateLimiter{
		storage: storage,
		config:  config,
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	count, err := rl.storage.Get(ctx, key)
	if err != nil {
		log.Printf("Key %s not found, starting new count", key)
		count = 0
	}

	log.Printf("Current count for key %s: %d, limit: %d", key, count, rl.config.RequestsPerSecond)

	// Check if the next increment would exceed the limit
	if count+1 > rl.config.RequestsPerSecond {
		log.Printf("Rate limit exceeded for key %s", key)
		return false, nil
	}

	newCount, err := rl.storage.Increment(ctx, key)
	if err != nil {
		return false, fmt.Errorf("failed to increment counter: %w", err)
	}

	log.Printf("New count after increment for key %s: %d", key, newCount)

	// Update TTL on every increment
	log.Printf("Setting expiration for key %s: %v", key, rl.config.BlockDuration)
	if err := rl.storage.Set(ctx, key, newCount, rl.config.BlockDuration); err != nil {
		return false, fmt.Errorf("failed to set expiration: %w", err)
	}

	return true, nil
}

func (rl *RateLimiter) Reset(ctx context.Context, key string) error {
	return rl.storage.Delete(ctx, key)
}
