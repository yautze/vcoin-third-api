package record

import (
	"context"
	"time"
	"vcoin-third-api/infra/redis"
)

//
const (
	db = 0
)

// Create -
func Create(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	// Set -
	if _, err := redis.Client.Do(ctx, "SELECT", db).Result(); err != nil {
		return err
	}

	_, err := redis.Client.Set(ctx, key, value, expire).Result()
	if err != nil {
		return err
	}

	return nil
}

// Get -
func Get(ctx context.Context, key string) (string, error) {
	if _, err := redis.Client.Do(ctx, "SELECT", db).Result(); err != nil {
		return "", err
	}

	res, err := redis.Client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}
