package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(host string, port int, password string, db int) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisStorage{client: client}, nil
}

func (r *RedisStorage) Increment(ctx context.Context, key string) (int64, error) {
	val, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key %s: %w", key, err)
	}
	log.Printf("Redis increment for key %s: %d", key, val)
	return val, nil
}

func (r *RedisStorage) Get(ctx context.Context, key string) (int64, error) {
	val, err := r.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		log.Printf("Key %s not found in Redis", key)
		return 0, fmt.Errorf("key not found")
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get key %s: %w", key, err)
	}
	log.Printf("Redis get for key %s: %d", key, val)
	return val, nil
}

func (r *RedisStorage) Set(ctx context.Context, key string, value int64, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	log.Printf("Redis set for key %s: %d with expiration %v", key, value, expiration)
	return nil
}

func (r *RedisStorage) Delete(ctx context.Context, key string) error {
	log.Printf("Deleting key: %s", key)
	return r.client.Del(ctx, key).Err()
}

func (r *RedisStorage) Clear(ctx context.Context) error {
	log.Printf("Clearing all keys from Redis")
	return r.client.FlushDB(ctx).Err()
}

func (r *RedisStorage) Close() error {
	return r.client.Close()
}
