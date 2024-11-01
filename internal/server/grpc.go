package server

import (
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
	server *grpc.Server
	port   string
	logger *slog.Logger
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

	return &GRPCServer{
		server: grpcServer,
		port:   port,
		logger: logger,
	}, nil
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
