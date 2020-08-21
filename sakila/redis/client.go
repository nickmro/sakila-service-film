package redis

import "github.com/go-redis/redis/v8"

// NewClient returns a new Redis client.
func NewClient(addr string, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
}
