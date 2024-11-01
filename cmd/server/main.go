package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
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

	// Create error channel for startup errors
	startupErrors := make(chan error, 3) // 3 for gRPC, gateway, and health check
	defer close(startupErrors)

	// Create WaitGroup for all servers
	var wg sync.WaitGroup
	wg.Add(3) // gRPC, gateway, and health check servers

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
		defer wg.Done()
		if err := healthServer.ListenAndServe(); err != http.ErrServerClosed {
			startupErrors <- fmt.Errorf("health check server error: %w", err)
		}
	}()

	// Start gRPC server
	go func() {
		defer wg.Done()
		if err := grpcServer.Start(); err != nil {
			startupErrors <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	// Start HTTP gateway server
	go func() {
		defer wg.Done()
		if err := gatewayServer.Start(); err != nil {
			startupErrors <- fmt.Errorf("gateway server error: %w", err)
		}
	}()

	// Wait for interrupt or startup error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Use separate goroutine to check for startup errors
	errChan := make(chan error, 1)
	go func() {
		select {
		case err := <-startupErrors:
			errChan <- err
		case <-quit:
			errChan <- nil
		}
	}()

	// Wait for shutdown signal or error
	if err := <-errChan; err != nil {
		logger.Error("Server startup failed", "error", err)
		os.Exit(1)
	}

	logger.Info("Initiating graceful shutdown")

	// Create root context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Create WaitGroup for shutdown operations
	var shutdownWg sync.WaitGroup
	shutdownWg.Add(3)

	// Stop health check server first
	go func() {
		defer shutdownWg.Done()
		if err := healthServer.Shutdown(ctx); err != nil {
			logger.Error("Failed to stop health check server", "error", err)
		}
	}()

	// Stop gateway server
	go func() {
		defer shutdownWg.Done()
		if err := gatewayServer.Stop(ctx); err != nil {
			logger.Error("Failed to stop gateway server", "error", err)
		}
	}()

	// Give in-flight requests time to reach gRPC server
	time.Sleep(shutdownGrace)

	// Stop gRPC server
	go func() {
		defer shutdownWg.Done()
		if err := grpcServer.Stop(ctx); err != nil {
			logger.Error("Failed to stop gRPC server", "error", err)
		}
	}()

	// Wait for all shutdown operations to complete or timeout
	shutdownComplete := make(chan struct{})
	go func() {
		shutdownWg.Wait()
		close(shutdownComplete)
	}()

	select {
	case <-ctx.Done():
		logger.Error("Shutdown timeout exceeded")
		os.Exit(1)
	case <-shutdownComplete:
		logger.Info("Graceful shutdown completed")
	}

	// Wait for all server goroutines to exit
	wg.Wait()
}
