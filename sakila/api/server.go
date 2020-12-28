package api

import (
	"fmt"
	"net/http"
	"sakila/sakila-film-service/sakila"
)

// Server is an http servier.
type Server struct {
	*http.Server
	Logger sakila.Logger
}

// NewServer returns a new server.
func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%s", addr),
			Handler: handler,
		},
	}
}
