package graphql

import (
	"net/http"

	"github.com/graphql-go/handler"
)

// NewHandler returns a new graphql http handler.
func NewHandler(s *Schema) http.Handler {
	return handler.New(&handler.Config{
		Schema:     s.Schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})
}
