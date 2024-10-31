package middleware

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// GRPCUnaryServerInterceptors returns a slice of unary server interceptors
func GRPCUnaryServerInterceptors(logger *slog.Logger) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		unaryLoggingInterceptor(logger),
		// Easy to add new middlewares here, for example:
		// rateLimitMiddleware(100),
		// metricMiddleware(),
		// authMiddleware(),
	}
}

// unaryLoggingInterceptor returns a new unary server interceptor that logs requests and responses
func unaryLoggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
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
