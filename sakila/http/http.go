// Package http provides http functions.
package http

import "net/http"

// ListenAndServe serves a handler over a given address.
func ListenAndServe(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}
