package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/ntayoane/go-rate-limited/internal/limiter"
)

func RateLimitMiddleware(ipLimiter, tokenLimiter *limiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
				ip = strings.Split(forwardedFor, ",")[0]
			}
			ip = strings.Split(ip, ":")[0] // Remove port number
			log.Printf("Request from IP: %s", ip)

			if token := r.Header.Get("API_KEY"); token != "" {
				log.Printf("Using token-based rate limiting for token: %s", token)
				allowed, err := tokenLimiter.Allow(r.Context(), token)
				if err != nil {
					log.Printf("Error in token rate limiting: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				if !allowed {
					log.Printf("Token rate limit exceeded for token: %s", token)
					http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
					return
				}
			} else {
				log.Printf("Using IP-based rate limiting for IP: %s", ip)
				allowed, err := ipLimiter.Allow(r.Context(), ip)
				if err != nil {
					log.Printf("Error in IP rate limiting: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				if !allowed {
					log.Printf("IP rate limit exceeded for IP: %s", ip)
					http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
