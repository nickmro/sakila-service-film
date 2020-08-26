package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Client is a Redis client.
type Client struct {
	*redis.Client
}

const pingTimeoutDuration = time.Second * 5

// NewClient returns a new Redis client.
func NewClient(addr string, password string) *Client {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	return &Client{Client: c}
}

// Status satisfies the health checker interface.
func (c *Client) Status() (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeoutDuration)
	defer cancel()

	return nil, c.Client.Ping(ctx).Err()
}
