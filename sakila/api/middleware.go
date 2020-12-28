package api

import (
	"encoding/json"
	"net/http"
	"sakila/sakila-film-service/sakila"
	"time"

	"github.com/go-chi/chi/middleware"
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

			start := time.Now()

			defer func() {
				log := &RequestLog{
					Path:         r.URL.Path,
					Protocol:     r.Proto,
					RequestID:    middleware.GetReqID(r.Context()),
					ResponseTime: time.Since(start).Milliseconds(),
					Size:         wrap.BytesWritten(),
					Status:       wrap.Status(),
				}

				b, err := json.Marshal(log)
				if err != nil {
					logger.Error(err)

					return
				}

				logger.Info(string(b))
			}()

			next.ServeHTTP(wrap, r)
		}

		return http.HandlerFunc(fn)
	}
}
