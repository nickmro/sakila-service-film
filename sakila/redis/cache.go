package redis

import (
	"context"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

// Cache is a redis cache.
type Cache struct {
	*cache.Cache
	client *redis.Client
}

const pingTimeoutDuration = time.Second * 10

// NewCache returns a new cache.
func NewCache(params *ClientParams) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address(params.Host, params.Port),
		Password: params.Password,
		DB:       params.DB,
	})

	status := client.Ping(context.Background())
	if err := status.Err(); err != nil {
		return nil, err
	}

	return &Cache{
		Cache: cache.New(&cache.Options{
			Redis:      client,
			LocalCache: cache.NewTinyLFU(10000, DefaultTTL),
		}),
		client: client,
	}, nil
}

// Status returns the client status.
func (cache *Cache) Status() (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeoutDuration)
	defer cancel()

	return nil, cache.client.Ping(ctx).Err()
}

// Close closes the cache client connection.
func (cache *Cache) Close() error {
	return cache.client.Close()
}
