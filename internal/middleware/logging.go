package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func LoggingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					slog.Error("Error in logging middleware", "err", err, "trace", debug.Stack())
				}
			}()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			slog.Info("Incoming request",
				"status", wrapped.status,
				"method", r.Method, "path",
				r.URL.EscapedPath(), "duration",
				time.Since(start))
		})
	}
}
