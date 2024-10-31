package server

import (
	"log/slog"
	"net"

	pb "github.com/fredrikaverpil/go-microservice/internal/proto/gen/go/gomicroservice/v1"
	"github.com/fredrikaverpil/go-microservice/internal/service"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server      *grpc.Server
	port        string
	logger      *slog.Logger
	userService *service.UserService
}

func NewGRPCServer(port string, logger *slog.Logger) *GRPCServer {
	grpcServer := grpc.NewServer()
	userService := service.NewUserService()

	// Register the user service
	pb.RegisterUserServiceServer(grpcServer, userService)

	return &GRPCServer{
		server:      grpcServer,
		port:        port,
		logger:      logger,
		userService: userService,
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
