package slidingWIndowLogs

import (
	"github.com/go-redis/redis/v9"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("UserID")
		uniqueRequestID := r.Header.Get("uniqueRequestID")

		if !RateLimitUsingSlidingLogs(userID, uniqueRequestID, 2, 5, redisClient) {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
