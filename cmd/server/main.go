package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bufbuild/protovalidate-go"
	"github.com/fredrikaverpil/go-microservice/internal/server"
)

const (
	grpcPort        = "50051"
	httpPort        = "8080"
	shutdownTimeout = 30 * time.Second
	shutdownGrace   = 2 * time.Second
)

func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Initialize proto validator
	validator, err := protovalidate.New()
	if err != nil {
		logger.Error("Failed to create proto validator", "error", err)
		os.Exit(1)
	}

	// Initialize gRPC server
	grpcServer, err := server.NewGRPCServer(grpcPort, logger, validator)
	if err != nil {
		logger.Error("Failed to create gRPC server", "error", err)
		os.Exit(1)
	}

	// Initialize HTTP gateway server
	gatewayServer, err := server.NewGatewayServer(httpPort, grpcPort, logger)
	if err != nil {
		logger.Error("Failed to create gateway server", "error", err)
		os.Exit(1)
	}

	// Add health check endpoint
	healthServer := &http.Server{
		Addr: ":8081",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !gatewayServer.HealthCheck() || !grpcServer.HealthCheck() {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			w.WriteHeader(http.StatusOK)
		}),
	}
	go func() {
		if err := healthServer.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("Health check server failed", "error", err)
		}
	}()

	// Start gRPC server
	go func() {
		if err := grpcServer.Start(); err != nil {
			logger.Error("Failed to start gRPC server", "error", err)
			os.Exit(1)
		}
	}()

	// Start HTTP gateway server
	go func() {
		if err := gatewayServer.Start(); err != nil {
			logger.Error("Failed to start gateway server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("Initiating graceful shutdown", "signal", sig)

	// Create root context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Stop health check server first
	if err := healthServer.Shutdown(ctx); err != nil {
		logger.Error("Failed to stop health check server", "error", err)
	}

	// First stop accepting new requests at gateway
	if err := gatewayServer.Stop(ctx); err != nil {
		logger.Error("Failed to stop gateway server", "error", err)
		os.Exit(1)
	}

	// Give in-flight requests time to reach gRPC server
	select {
	case <-time.After(shutdownGrace):
	case <-ctx.Done():
		logger.Error("Shutdown grace period exceeded")
	}

	// Now stop the gRPC server
	if err := grpcServer.Stop(ctx); err != nil {
		logger.Error("Failed to stop gRPC server", "error", err)
		os.Exit(1)
	}

	// Wait for context to ensure we don't exceed total shutdown timeout
	select {
	case <-ctx.Done():
		logger.Error("Shutdown timeout exceeded")
		os.Exit(1)
	default:
		logger.Info("Graceful shutdown completed")
	}
}
