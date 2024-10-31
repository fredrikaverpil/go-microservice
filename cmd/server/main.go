package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/fredrikaverpil/go-microservice/internal/server"
)

func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	grpcServer := server.NewGRPCServer("50051", logger)

	// Start server in a goroutine
	go func() {
		if err := grpcServer.Start(); err != nil {
			logger.Error("Failed to start gRPC server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// Gracefully stop the server
	grpcServer.Stop()
}
