package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// HTTPMiddleware is a type for HTTP middleware functions
type HTTPMiddleware func(http.Handler) http.Handler

// Chain applies a list of middlewares to a handler
func Chain(middlewares ...HTTPMiddleware) HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// HTTPServerMiddlewares returns a list of HTTP server middlewares
func HTTPServerMiddlewares(logger *slog.Logger) []HTTPMiddleware {
	return []HTTPMiddleware{
		loggingMiddleware(logger),
	}
}

// loggingMiddleware returns a new middleware that logs requests and responses
func loggingMiddleware(logger *slog.Logger) HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Log request
			logger.Info("received request",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
			)

			// Create custom response writer to capture status code
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			// Handle request
			next.ServeHTTP(rw, r)

			// Log response
			logger.Info("request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.status,
				"duration", time.Since(start),
			)
		})
	}
}
