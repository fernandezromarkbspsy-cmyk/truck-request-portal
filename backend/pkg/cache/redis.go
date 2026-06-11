package cache

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() error {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		return fmt.Errorf("REDIS_URL environment variable is not set")
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return fmt.Errorf("unable to parse redis url: %v", err)
	}

	RedisClient = redis.NewClient(opts)
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("unable to ping redis: %v", err)
	}

	return nil
}
