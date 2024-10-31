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

// WithHTTPMiddlewares chains multiple HTTP middlewares together
func WithHTTPMiddlewares(handler http.Handler, middlewares ...HTTPMiddleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// HTTPServerMiddlewares returns a list of HTTP server middlewares
func HTTPServerMiddlewares(logger *slog.Logger) []HTTPMiddleware {
	return []HTTPMiddleware{
		loggingMiddleware(logger),
		// Easy to add new middlewares here, for example:
		// rateLimitMiddleware(100),
		// metricMiddleware(),
		// authMiddleware(),
	}
}

// loggingMiddleware handles request/response logging
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

// Example of how to add a new middleware:
/*
func rateLimitMiddleware(rps int) HTTPMiddleware {
	limiter := rate.NewLimiter(rate.Limit(rps), rps)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
