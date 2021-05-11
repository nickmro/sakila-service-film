// Package redis handles operations on the Redis cache.
package redis

import (
	"crypto/sha1"
	"fmt"
	"time"
)

// DefaultTTL is the default cache TTL.
const DefaultTTL = time.Minute * 5

func hashedKey(key string) string {
	h := sha1.New()

	if _, err := h.Write([]byte(key)); err != nil {
		return ""
	}

	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}
