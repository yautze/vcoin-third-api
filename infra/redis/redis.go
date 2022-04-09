package redis

import (
	"context"

	goredis "github.com/go-redis/redis/v8"
)

// Client -
var Client *goredis.Client

// New -
func New() error {
	ctx := context.Background()

	opts := goredis.Options{
		Addr:         "localhost:6379",
		Password:     "",
		PoolSize:     20,
		MaxRetries:   3,
		MinIdleConns: 5,
	}

	client := goredis.NewClient(&opts)
	if _, err := client.Ping(ctx).Result(); err != nil {
		return err
	}

	Client = client

	return nil
}
