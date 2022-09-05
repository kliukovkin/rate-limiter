package tokenBucket

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v9"
	"time"
)

type bucket struct {
	Value      int64 `json:"value"`
	LastUpdate int64 `json:"lastUpdate"`
}

func (i bucket) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func RateLimitUsingTokenBucket(userID string, intervalInSeconds int64, maximumRequests int64, redisClient *redis.Client) bool {
	ctx := context.Background()
	key := "user:" + userID + ":rate-limit"
	value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		b := bucket{
			Value:      maximumRequests,
			LastUpdate: time.Now().Unix(),
		}
		redisClient.Set(ctx, key, b, 0)
		return true
	}
	var b bucket
	err = json.Unmarshal([]byte(value), &b)

	if time.Now().Unix()-b.LastUpdate >= intervalInSeconds {
		b.Value = maximumRequests
		b.LastUpdate = time.Now().Unix()
		redisClient.Set(ctx, key, b, 0)
	} else {
		if b.Value <= 0 { // request left is 0 or < 0
			// drop request
			return false
		}
	}

	b.Value--
	redisClient.Set(ctx, key, b, 0)

	return true
	//TODO add refill logic
	//TODO add race condition handler for redis(incr/decr)
}

func RateLimitUsingLeakyBucket(
	userID string,
	uniqueRequestID string,
	intervalInSeconds int64,
	maximumRequests int64,
	redisClient *redis.Client) bool {
	ctx := context.Background()
	// userID can be apikey, location, ip
	requestCount := redisClient.LLen(ctx, userID).Val()

	if requestCount >= maximumRequests {
		// drop request
		return false
	}

	// add request id to the end of request queue
	redisClient.RPush(ctx, userID, uniqueRequestID)

	return true
}
