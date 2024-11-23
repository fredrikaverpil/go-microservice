package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// HTTP logging middleware.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

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

// gRPC logging middleware.
func unaryLoggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Log request
		logger.Info("received request",
			"method", info.FullMethod,
			"request", req,
		)

		// Handle request
		resp, err := handler(ctx, req)

		// Log response
		if err != nil {
			st, _ := status.FromError(err)
			logger.Error("request failed",
				"method", info.FullMethod,
				"duration", time.Since(start),
				"code", st.Code(),
				"error", err,
			)
		} else {
			logger.Info("request succeeded",
				"method", info.FullMethod,
				"duration", time.Since(start),
				"response", resp,
			)
		}

		return resp, err
	}
}
