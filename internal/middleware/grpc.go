package middleware

import (
	"log/slog"

	"google.golang.org/grpc"
)

// GRPCUnaryServerInterceptors returns a slice of unary server interceptors
func GRPCUnaryServerInterceptors(logger *slog.Logger) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		// circuitBreakerUnaryInterceptor(),
		unaryLoggingInterceptor(logger),
		// rateLimitInterceptor(100),
		// metricInterceptor(),
		// authInterceptor(),
	}
}
