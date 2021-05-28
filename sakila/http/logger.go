package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/nickmro/sakila-service-film/sakila"
)

// RequestLog contains request information.
type RequestLog struct {
	Path         string `json:"path"`
	Protocol     string `json:"protocol"`
	RequestID    string `json:"request_id"`
	ResponseTime int64  `json:"response_time"`
	Size         int    `json:"size"`
	Status       int    `json:"status"`
}

// RequestLogger returns a request logger middleware.
func RequestLogger(logger sakila.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			wrap := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			reqID := middleware.GetReqID(r.Context())

			start := time.Now()

			defer func() {
				log := &RequestLog{
					Path:         r.URL.Path,
					Protocol:     r.Proto,
					RequestID:    reqID,
					ResponseTime: time.Since(start).Milliseconds(),
					Size:         wrap.BytesWritten(),
					Status:       wrap.Status(),
				}

				if b, err := json.Marshal(log); err == nil {
					logger.Info(string(b))
				}
			}()

			next.ServeHTTP(wrap, r)
		}

		return http.HandlerFunc(fn)
	}
}
