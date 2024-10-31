package server

import (
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	server *grpc.Server
	port   string
	logger *slog.Logger
}

func NewGRPCServer(port string, logger *slog.Logger) *GRPCServer {
	return &GRPCServer{
		server: grpc.NewServer(),
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
