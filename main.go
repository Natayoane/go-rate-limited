package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ntayoane/go-rate-limited/internal/config"
	"github.com/ntayoane/go-rate-limited/internal/limiter"
	"github.com/ntayoane/go-rate-limited/internal/middleware"
	"github.com/ntayoane/go-rate-limited/internal/storage"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	redisStorage, err := storage.NewRedisStorage(
		cfg.RedisHost,
		cfg.RedisPort,
		cfg.RedisPassword,
		cfg.RedisDB,
	)
	if err != nil {
		log.Fatalf("Failed to initialize Redis storage: %v", err)
	}
	defer redisStorage.Close()

	ipLimiter := limiter.NewRateLimiter(redisStorage, limiter.Config{
		RequestsPerSecond: cfg.IPRequestsPerSecond,
		BlockDuration:     cfg.IPBlockDuration,
	})

	tokenLimiter := limiter.NewRateLimiter(redisStorage, limiter.Config{
		RequestsPerSecond: cfg.TokenRequestsPerSecond,
		BlockDuration:     cfg.TokenBlockDuration,
	})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: middleware.RateLimitMiddleware(ipLimiter, tokenLimiter)(handler),
	}

	go func() {
		log.Printf("Server starting on port %d", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")
}
