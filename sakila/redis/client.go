package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Client is a redis client.
type Client struct {
	*redis.Client
}

// ClientParams are Redis client parameters.
type ClientParams struct {
	Host     string
	Port     int
	Password string
	DB       int
}

const pingTimeoutDuration = time.Second * 10

func address(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

// NewClient returns a new client.
func NewClient(params *ClientParams) *Client {
	return &Client{
		Client: redis.NewClient(&redis.Options{
			Addr:     address(params.Host, params.Port),
			Password: params.Password,
			DB:       params.DB,
		}),
	}
}

// Status returns the client status.
func (c *Client) Status() (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeoutDuration)
	defer cancel()

	return nil, c.Client.Ping(ctx).Err()
}
