package middleware

import (
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

			if token := r.Header.Get("API_KEY"); token != "" {
				allowed, err := tokenLimiter.Allow(r.Context(), token)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				if !allowed {
					http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
					return
				}
			} else {
				allowed, err := ipLimiter.Allow(r.Context(), ip)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				if !allowed {
					http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
