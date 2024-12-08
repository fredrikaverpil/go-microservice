package integration

import (
	"context"
	"log/slog"
	"os"

	"github.com/bufbuild/protovalidate-go"
	"github.com/fredrikaverpil/go-microservice/internal/core/service"
	"github.com/fredrikaverpil/go-microservice/internal/inbound/handler/grpc/gomicroservice"
	"github.com/fredrikaverpil/go-microservice/internal/outbound/db"
)

type fixture struct {
	userHandler *gomicroservice.GRPCHandler
	ctx         context.Context
}

func newFixture() *fixture {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	validator, err := protovalidate.New()
	if err != nil {
		logger.Error("Failed to create proto validator", "error", err)
		os.Exit(1)
	}

	userRepo := db.NewMemoryRepository(logger)
	userService := service.NewUserService(logger, userRepo)
	userHandler := gomicroservice.NewGRPCHandler(userService, validator)

	return &fixture{
		userHandler: userHandler,
		ctx:         context.Background(),
	}
}
