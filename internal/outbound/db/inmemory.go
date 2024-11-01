package db

import (
	"context"
	"log/slog"
	"sync"

	"github.com/fredrikaverpil/go-microservice/internal/core/domain"
	"github.com/fredrikaverpil/go-microservice/internal/core/port"
)

type MemoryRepository struct {
	users  map[string]*domain.User
	mutex  sync.RWMutex
	logger *slog.Logger
}

func NewMemoryRepository(logger *slog.Logger) port.UserRepository {
	return &MemoryRepository{
		users:  make(map[string]*domain.User),
		logger: logger,
	}
}

func (r *MemoryRepository) CreateUser(
	ctx context.Context,
	user *domain.User,
) (*domain.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return nil, domain.NewErrorNotFound("user not found", nil)
}

func (r *MemoryRepository) GetUser(ctx context.Context, name string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[name]
	if !exists {
		return nil, domain.NewErrorNotFound("user not found", nil)
	}
	return user, nil
}

func (r *MemoryRepository) ListUsers(
	ctx context.Context,
	pageSize int32,
	pageToken string,
) ([]*domain.User, string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return nil, "", domain.NewErrorNotFound("user not found", nil)
}

func (r *MemoryRepository) UpdateUser(
	ctx context.Context,
	user *domain.User,
) (*domain.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return nil, domain.NewErrorNotFound("user not found", nil)
}

func (r *MemoryRepository) DeleteUser(ctx context.Context, name string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return domain.NewErrorNotFound("user not found", nil)
}
