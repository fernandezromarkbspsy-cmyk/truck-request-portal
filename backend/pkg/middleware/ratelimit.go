package middleware

import (
	"context"
	"net/http"
	"time"
	"truck-request-portal/pkg/cache"
)

// RateLimit prevents API abuse by limiting requests per IP/User.
func RateLimit(maxRequests int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Use IP address as the key (in production, use UserID for authenticated routes)
			key := "ratelimit:" + r.RemoteAddr

			// Increment counter in Redis
			_, err := cache.RedisClient.Incr(context.Background(), key).Result()
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Set expiration on first request
			ttl, _ := cache.RedisClient.TTL(context.Background(), key).Result()
			if ttl < 0 {
				cache.RedisClient.Expire(context.Background(), key, window)
			}

			// Check if limit exceeded
			count, _ := cache.RedisClient.Get(context.Background(), key).Int()
			if count > maxRequests {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
