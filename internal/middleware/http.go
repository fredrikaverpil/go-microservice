package middleware

import (
	"log/slog"
	"net/http"
)

// HTTPMiddleware is a type for HTTP middleware functions.
type HTTPMiddleware func(http.Handler) http.Handler

// WithHTTPMiddlewares chains multiple HTTP middlewares together.
func WithHTTPMiddlewares(handler http.Handler, middlewares ...HTTPMiddleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// HTTPServerMiddlewares returns a list of HTTP server middlewares.
func HTTPServerMiddlewares(logger *slog.Logger) []HTTPMiddleware {
	return []HTTPMiddleware{
		// circuitBreakerMiddleware(),
		loggingMiddleware(logger),
		// rateLimitMiddleware(100),
		// metricMiddleware(),
		// authMiddleware(),
	}
}
