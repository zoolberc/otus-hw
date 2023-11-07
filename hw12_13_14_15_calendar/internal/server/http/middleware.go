package internalhttp

import (
	"net/http"
	"time"

	"golang.org/x/exp/slog"
)

func loggingMiddleware(next http.HandlerFunc, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entry := log.With(
			slog.String("clientIp", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("userAgent", r.UserAgent()),
			slog.String("path", r.URL.Path),
			slog.String("httpVersion", r.Proto),
		)

		start := time.Now()
		defer func() {
			entry.Info("request completed",
				slog.String("status", "??"),
				slog.Int64("latency", time.Since(start).Milliseconds()),
			)
		}()

		next.ServeHTTP(w, r)
	}
}
