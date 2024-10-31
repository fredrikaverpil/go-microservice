package server

import (
	"context"
	"log/slog"
	"net/http"

	gomicroservicev1 "github.com/fredrikaverpil/go-microservice/internal/proto/gen/go/gomicroservice/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GatewayServer struct {
	server *http.Server
	logger *slog.Logger
}

func NewGatewayServer(port string, grpcPort string, logger *slog.Logger) (*GatewayServer, error) {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	endpoint := "localhost:" + grpcPort

	if err := gomicroservicev1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	return &GatewayServer{
		server: server,
		logger: logger,
	}, nil
}

func (s *GatewayServer) Start() error {
	s.logger.Info("HTTP gateway server listening", "port", s.server.Addr)
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *GatewayServer) Stop(ctx context.Context) error {
	s.logger.Info("HTTP gateway server stopping")
	return s.server.Shutdown(ctx)
}
