package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// ListenAndServe serves the given handler on the given port.
func ListenAndServe(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}

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

	return http.NewRequest(method, url, r)
}
