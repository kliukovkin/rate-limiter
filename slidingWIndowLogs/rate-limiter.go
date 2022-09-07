package slidingWIndowLogs

import (
	"context"
	"github.com/go-redis/redis/v9"
	"strconv"
	"time"
)

// RateLimitUsingSlidingLogs .
func RateLimitUsingSlidingLogs(userID string, uniqueRequestID string, intervalInSeconds int64, maximumRequests int64, redisClient *redis.Client) bool {
	ctx := context.Background()
	// userID can be apikey, location, ip
	currentTime := strconv.FormatInt(time.Now().Unix(), 10)
	lastWindowTime := strconv.FormatInt(time.Now().Unix()-intervalInSeconds, 10)
	// get current window count
	requestCount := redisClient.ZCount(ctx, userID, lastWindowTime, currentTime).Val()
	if requestCount >= maximumRequests {
		// drop request
		return false
	}

	// add request id to last window
	redisClient.ZAdd(ctx, userID, redis.Z{Score: float64(time.Now().Unix()), Member: uniqueRequestID}) //TODO нахуя тут uniqueRequestID???
	// add unique request id to userID set with score as current time

	// handle request
	return true
	//TODO remove all expired request ids at regular intervals using ZRemRangeByScore from -inf to last window time
}
