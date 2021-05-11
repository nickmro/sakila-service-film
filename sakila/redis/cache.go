package redis

import (
	"context"

	"github.com/go-redis/cache/v8"
)

// NewCache returns a new cache.
func NewCache(client *Client) (*cache.Cache, error) {
	status := client.Ping(context.Background())
	if err := status.Err(); err != nil {
		return nil, err
	}

	return cache.New(&cache.Options{
		Redis:      client,
		LocalCache: cache.NewTinyLFU(10000, DefaultTTL),
	}), nil
}
