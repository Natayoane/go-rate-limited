package limiter

import (
	"context"
	"fmt"
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
	count, err := rl.storage.Increment(ctx, key)
	if err != nil {
		return false, fmt.Errorf("failed to increment counter: %w", err)
	}

	if count == 1 {
		if err := rl.storage.Set(ctx, key, 1, rl.config.BlockDuration); err != nil {
			return false, fmt.Errorf("failed to set expiration: %w", err)
		}
	}

	return count <= rl.config.RequestsPerSecond, nil
}

func (rl *RateLimiter) Reset(ctx context.Context, key string) error {
	return rl.storage.Delete(ctx, key)
}
