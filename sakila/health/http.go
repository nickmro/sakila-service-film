package health

import (
	"net/http"

	"github.com/InVisionApp/go-health/v2/handlers"
)

// NewHandler returns a new healthcheck handler.
func NewHandler(checker *Checker) http.HandlerFunc {
	return handlers.NewJSONHandlerFunc(checker, nil)
}
