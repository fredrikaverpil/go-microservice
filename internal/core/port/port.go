package port

import (
	"context"

	"github.com/fredrikaverpil/go-microservice/internal/core/domain"
)

type UserService interface { //nolint: iface // UserService/UserRepository equal today but may diverge in the future.
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUser(ctx context.Context, name string) (*domain.User, error)
	ListUsers(ctx context.Context, pageSize int32, pageToken string) ([]*domain.User, string, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, name string) error
}

type UserRepository interface { //nolint: iface // UserService/UserRepository equal today but may diverge in the future.
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUser(ctx context.Context, name string) (*domain.User, error)
	ListUsers(ctx context.Context, pageSize int32, pageToken string) ([]*domain.User, string, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, name string) error
}
