package service

import (
	"context"
	"log/slog"

	"github.com/fredrikaverpil/go-microservice/internal/domain"
	"github.com/fredrikaverpil/go-microservice/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	logger *slog.Logger
}

func NewUserService(logger *slog.Logger) ports.UserService {
	return &UserService{
		logger: logger,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	err := status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
	s.logger.Error("failed to create user",
		"error", err,
		"user", user.Name,
	)
	return nil, err
}

func (s *UserService) GetUser(ctx context.Context, name string) (*domain.User, error) {
	err := status.Errorf(codes.Unimplemented, "method GetUser not implemented")
	s.logger.Error("failed to get user",
		"error", err,
		"name", name,
	)
	return nil, err
}

func (s *UserService) ListUsers(
	ctx context.Context,
	pageSize int32,
	pageToken string,
) ([]*domain.User, string, error) {
	err := status.Errorf(codes.Unimplemented, "method ListUsers not implemented")
	s.logger.Error("failed to list users",
		"error", err,
		"pageSize", pageSize,
		"pageToken", pageToken,
	)
	return nil, "", err
}

func (s *UserService) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	err := status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
	s.logger.Error("failed to update user",
		"error", err,
		"user", user.Name,
	)
	return nil, err
}

func (s *UserService) DeleteUser(ctx context.Context, name string) error {
	err := status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
	s.logger.Error("failed to delete user",
		"error", err,
		"name", name,
	)
	return err
}
