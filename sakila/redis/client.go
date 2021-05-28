package redis

import "fmt"

// ClientParams are Redis client parameters.
type ClientParams struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func address(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}
