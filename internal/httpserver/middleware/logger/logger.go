// Package logger provides HTTP middleware for request logging.
package logger

import (
	"net/http"
	"time"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
)

// New creates a new HTTP middleware for logging requests.
func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			// Collect initial request information
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			// Wrap http.ResponseWriter to capture response details
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Record request start time to calculate processing duration
			t1 := time.Now()

			// Log will be written in defer after request is processed
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			// Pass control to the next handler in the middleware chain
			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
