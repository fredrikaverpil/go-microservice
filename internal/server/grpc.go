package server

import (
	"context"
	"log/slog"
	"net"

	"github.com/bufbuild/protovalidate-go"
	"github.com/fredrikaverpil/go-microservice/internal/config"
	"github.com/fredrikaverpil/go-microservice/internal/core/service"
	pb "github.com/fredrikaverpil/go-microservice/internal/inbound/handler/grpc/gen/go/gomicroservice/v1"
	"github.com/fredrikaverpil/go-microservice/internal/inbound/handler/grpc/gomicroservice"
	"github.com/fredrikaverpil/go-microservice/internal/middleware"
	"github.com/fredrikaverpil/go-microservice/internal/outbound/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	server   *grpc.Server
	port     string
	logger   *slog.Logger
	listener net.Listener
	ready    bool
}

func NewGRPCServer(
	port string,
	logger *slog.Logger,
	validator *protovalidate.Validator,
) (*GRPCServer, error) {
	// Create server with interceptors
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middleware.GRPCUnaryServerInterceptors(logger)...),
	)

	// Create repository and service
	userRepo := db.NewMemoryRepository(logger)
	userService := service.NewUserService(logger, userRepo)
	userHandler := gomicroservice.NewGRPCHandler(userService, validator)

	// Register handler
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	// Enable reflection in development
	if config.IsDevelopment() {
		reflection.Register(grpcServer)
		logger.Info("gRPC reflection enabled")
	}

	// Create listener during initialization
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}

	return &GRPCServer{
		server:   grpcServer,
		port:     port,
		logger:   logger,
		listener: lis,
		ready:    false,
	}, nil
}

func (s *GRPCServer) Start() error {
	s.logger.Info("gRPC server listening", "port", s.port)
	s.ready = true
	if err := s.server.Serve(s.listener); err != nil {
		s.ready = false
		return err
	}
	return nil
}

func (s *GRPCServer) Stop(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		s.ready = false
		s.server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.server.Stop()
		return ctx.Err()
	case <-stopped:
		s.logger.Info("gRPC server stopped gracefully")
		return nil
	}
}

func (s *GRPCServer) IsReady() bool {
	return s.ready
}

func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}
