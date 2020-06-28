package tools

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

// SetValue is a function that adds a key value pair to a redis database
func SetValue(redisClient *redis.Client, key string, value string, expiry time.Duration) error {
	ctx := context.Background()
	err := redisClient.Set(ctx, key, value, expiry).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetValue is a function that searchs for a value of a provided key on a redis database
func GetValue(redisClient *redis.Client, key string) (string, error) {
	ctx := context.Background()
	value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

// RemoveValues is a function that removes a key value pair from a redis database
func RemoveValues(redisClient *redis.Client, key ...string) {
	ctx := context.Background()
	redisClient.Del(ctx, key...)
}
