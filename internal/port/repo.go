package port

import (
	"context"

	"github.com/fredrikaverpil/go-microservice/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUser(ctx context.Context, name string) (*domain.User, error)
	ListUsers(ctx context.Context, pageSize int32, pageToken string) ([]*domain.User, string, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, name string) error
}
