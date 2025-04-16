package server

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/fredrikaverpil/go-microservice/internal/config"
	gomicroservicev1 "github.com/fredrikaverpil/go-microservice/internal/gen/gomicroservice/v1"
	"github.com/fredrikaverpil/go-microservice/internal/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	readHeaderTimeout = 5 * time.Second
)

type GatewayServer struct {
	server *http.Server
	logger *slog.Logger
	state  State
	ready  bool
}

func NewGatewayServer(
	port, grpcPort string,
	logger *slog.Logger,
) (*GatewayServer, error) {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	// Create client connection to gRPC server
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	endpoint := "localhost:" + grpcPort

	if err := gomicroservicev1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return nil, err
	}

	swaggerHandler := SwaggerHandler(logger)

	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Route Swagger UI and OpenAPI spec requests
		if config.IsDevelopment() && (strings.HasPrefix(r.URL.Path, "/docs") || r.URL.Path == "/api/openapi.yaml") {
			swaggerHandler.ServeHTTP(w, r)
			return
		}

		// All other paths go to the gRPC-gateway
		mux.ServeHTTP(w, r)
	})

	// Wrap mux with middlewares
	handler := middleware.WithHTTPMiddlewares(mainHandler, middleware.HTTPServerMiddlewares(logger)...)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	return &GatewayServer{
		server: server,
		logger: logger,
		state:  StateStarting,
		ready:  false,
	}, nil
}

func (s *GatewayServer) Start() error {
	s.logger.Info("HTTP gateway server listening", "port", s.server.Addr)
	s.ready = true
	s.state = StateRunning
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		s.ready = false
		s.state = StateStopped
		return err
	}
	return nil
}

func (s *GatewayServer) Stop(ctx context.Context) error {
	s.state = StateShuttingDown
	s.ready = false
	s.logger.InfoContext(ctx, "HTTP gateway server stopping")
	err := s.server.Shutdown(ctx)
	s.state = StateStopped
	return err
}

func (s *GatewayServer) IsReady() bool {
	return s.ready
}

func (s *GatewayServer) HealthCheck() bool {
	return s.state == StateRunning && s.ready
}

func (s *GatewayServer) State() State {
	return s.state
}
