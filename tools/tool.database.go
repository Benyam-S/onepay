package tools

import (
	"context"
	"errors"
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
	// should refine and analyze the key

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

// AnalyzeKeyValuePair is a function that cross checks the given key value pair with the stored one
func AnalyzeKeyValuePair(redisClient *redis.Client, key, value string) error {

	// Checking for empty values
	if len(key) == 0 || len(value) == 0 {
		return errors.New("empty key value pair used")
	}

	// Retriving value from redis store
	storedValue, err := GetValue(redisClient, key)
	if err != nil {
		return errors.New("value not found")
	}

	// Checking if the provided value matches the value from the database
	if storedValue != value {
		return errors.New("value does not match")
	}

	// Removing key value pair from the redis store
	RemoveValues(redisClient, key)

	return nil
}
