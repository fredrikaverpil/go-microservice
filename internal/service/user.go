package service

import (
	"context"
	"log/slog"

	"github.com/fredrikaverpil/go-microservice/internal/domain"
	"github.com/fredrikaverpil/go-microservice/internal/port"
)

type UserService struct {
	logger *slog.Logger
	repo   port.UserRepository
}

func NewUserService(logger *slog.Logger, repo port.UserRepository) port.UserService {
	return &UserService{
		logger: logger,
		repo:   repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		s.logger.Error("failed to create user",
			"error", err,
			"user", user.Name,
		)
		return nil, err // Propagate the custom error
	}
	return createdUser, nil
}

func (s *UserService) GetUser(ctx context.Context, name string) (*domain.User, error) {
	user, err := s.repo.GetUser(ctx, name)
	if err != nil {
		s.logger.Error("failed to get user",
			"error", err,
			"name", name,
		)
		return nil, err // Propagate the custom error
	}
	return user, nil
}

func (s *UserService) ListUsers(
	ctx context.Context,
	pageSize int32,
	pageToken string,
) ([]*domain.User, string, error) {
	users, nextToken, err := s.repo.ListUsers(ctx, pageSize, pageToken)
	if err != nil {
		s.logger.Error("failed to list users",
			"error", err,
			"pageSize", pageSize,
			"pageToken", pageToken,
		)
		return nil, "", err // Propagate the custom error
	}
	return users, nextToken, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	updatedUser, err := s.repo.UpdateUser(ctx, user)
	if err != nil {
		s.logger.Error("failed to update user",
			"error", err,
			"user", user.Name,
		)
		return nil, err // Propagate the custom error
	}
	return updatedUser, nil
}

func (s *UserService) DeleteUser(ctx context.Context, name string) error {
	if err := s.repo.DeleteUser(ctx, name); err != nil {
		s.logger.Error("failed to delete user",
			"error", err,
			"name", name,
		)
		return err // Propagate the custom error
	}
	return nil
}
