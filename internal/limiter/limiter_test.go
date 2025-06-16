package limiter

import (
	"context"
	"testing"
	"time"
)

type MockStorage struct {
	counts map[string]int64
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		counts: make(map[string]int64),
	}
}

func (m *MockStorage) Increment(ctx context.Context, key string) (int64, error) {
	m.counts[key]++
	return m.counts[key], nil
}

func (m *MockStorage) Get(ctx context.Context, key string) (int64, error) {
	return m.counts[key], nil
}

func (m *MockStorage) Set(ctx context.Context, key string, value int64, expiration time.Duration) error {
	m.counts[key] = value
	return nil
}

func (m *MockStorage) Delete(ctx context.Context, key string) error {
	delete(m.counts, key)
	return nil
}

func (m *MockStorage) Close() error {
	return nil
}

func TestRateLimiter(t *testing.T) {
	storage := NewMockStorage()
	config := Config{
		RequestsPerSecond: 5,
		BlockDuration:     time.Second * 1,
	}
	limiter := NewRateLimiter(storage, config)

	tests := []struct {
		name     string
		key      string
		requests int
		want     bool
	}{
		{
			name:     "Under limit",
			key:      "test-key-1",
			requests: 3,
			want:     true,
		},
		{
			name:     "At limit",
			key:      "test-key-2",
			requests: 5,
			want:     true,
		},
		{
			name:     "Over limit",
			key:      "test-key-3",
			requests: 6,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make the specified number of requests
			var lastResult bool
			for i := 0; i < tt.requests; i++ {
				allowed, err := limiter.Allow(context.Background(), tt.key)
				if err != nil {
					t.Errorf("RateLimiter.Allow() error = %v", err)
					return
				}
				lastResult = allowed
			}

			if lastResult != tt.want {
				t.Errorf("RateLimiter.Allow() = %v, want %v", lastResult, tt.want)
			}
		})
	}
}
