package http

import (
	"net/http"
	"sakila/sakila-film-service/sakila/graphql"

	"github.com/graphql-go/handler"
)

// NewGraphQLHandler returns a new graphql handler.
func NewGraphQLHandler(s *graphql.Schema) http.Handler {
	return handler.New(&handler.Config{
		Schema:     s.Schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})
}
