package server

import (
	"log/slog"
	"net"

	"github.com/fredrikaverpil/go-microservice/internal/config"
	"github.com/fredrikaverpil/go-microservice/internal/handler"
	"github.com/fredrikaverpil/go-microservice/internal/middleware"
	pb "github.com/fredrikaverpil/go-microservice/internal/proto/gen/go/gomicroservice/v1"
	"github.com/fredrikaverpil/go-microservice/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	server *grpc.Server
	port   string
	logger *slog.Logger
}

func NewGRPCServer(port string, logger *slog.Logger) *GRPCServer {
	// Create server with logging interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryServerLoggingInterceptor(logger)),
	)

	// Create service and handler with logger
	userService := service.NewUserService(logger)
	userHandler := handler.NewGRPCHandler(userService)

	// Register handler
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	// Enable reflection in development
	if config.IsDevelopment() {
		reflection.Register(grpcServer)
		logger.Info("gRPC reflection enabled")
	}

	return &GRPCServer{
		server: grpcServer,
		port:   port,
		logger: logger,
	}
}

func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return err
	}

	s.logger.Info("gRPC server listening", "port", s.port)
	if err := s.server.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
	s.logger.Info("gRPC server stopped")
}

func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}
