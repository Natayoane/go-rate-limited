package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	IPRequestsPerSecond    int64
	IPBlockDuration        time.Duration
	TokenRequestsPerSecond int64
	TokenBlockDuration     time.Duration
	RedisHost              string
	RedisPort              int
	RedisPassword          string
	RedisDB                int
	ServerPort             int
}

func LoadConfig() (*Config, error) {
	config := &Config{}

	ipRequests, err := strconv.ParseInt(getEnv("IP_REQUESTS_PER_SECOND", "5"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid IP_REQUESTS_PER_SECOND: %w", err)
	}
	config.IPRequestsPerSecond = ipRequests

	ipBlockDuration, err := strconv.ParseInt(getEnv("IP_BLOCK_DURATION_SECONDS", "300"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid IP_BLOCK_DURATION_SECONDS: %w", err)
	}
	config.IPBlockDuration = time.Duration(ipBlockDuration) * time.Second

	tokenRequests, err := strconv.ParseInt(getEnv("DEFAULT_TOKEN_REQUESTS_PER_SECOND", "10"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid DEFAULT_TOKEN_REQUESTS_PER_SECOND: %w", err)
	}
	config.TokenRequestsPerSecond = tokenRequests

	tokenBlockDuration, err := strconv.ParseInt(getEnv("DEFAULT_TOKEN_BLOCK_DURATION_SECONDS", "300"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid DEFAULT_TOKEN_BLOCK_DURATION_SECONDS: %w", err)
	}
	config.TokenBlockDuration = time.Duration(tokenBlockDuration) * time.Second

	config.RedisHost = getEnv("REDIS_HOST", "localhost")

	redisPort, err := strconv.Atoi(getEnv("REDIS_PORT", "6379"))
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_PORT: %w", err)
	}
	config.RedisPort = redisPort

	config.RedisPassword = getEnv("REDIS_PASSWORD", "")

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
	}
	config.RedisDB = redisDB

	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}
	config.ServerPort = serverPort

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
