// Package api provides handlers for the API server.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// NewRequest returns a new request.
func NewRequest(method string, url string, body interface{}) (*http.Request, error) {
	var r io.Reader

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		r = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return nil, err
	}

	return req.WithContext(context.TODO()), nil
}

// NewRequestWithContext returns a new request with a context.
func NewRequestWithContext(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	req, err := NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	return req.WithContext(ctx), nil
}
