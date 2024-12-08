package db

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

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
	_ context.Context,
	user *domain.User,
) (*domain.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Validate required fields
	if user.DisplayName == "" {
		return nil, domain.NewErrorInvalidInput("display_name is required", nil)
	}
	if user.Email == "" {
		return nil, domain.NewErrorInvalidInput("email is required", nil)
	}

	// Check if user already exists
	if _, exists := r.users[user.Name]; exists {
		return nil, domain.NewErrorAlreadyExists(
			fmt.Sprintf("user already exists: %s", user.Name),
			nil,
		)
	}

	// Create new user with timestamps
	now := time.Now().UTC()
	newUser := &domain.User{
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		CreateTime:  now,
		UpdateTime:  now,
	}

	// Store user
	r.users[newUser.Name] = newUser

	// Return a copy to prevent external modifications
	copyUser, err := newUser.Copy()
	if err != nil {
		return nil, domain.NewErrorInternal("failed to copy user", err)
	}
	return copyUser, nil
}

func (r *MemoryRepository) GetUser(_ context.Context, name string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[name]
	if !exists || !user.DeleteTime.IsZero() {
		return nil, domain.NewErrorNotFound("user not found", nil)
	}
	return user, nil
}

func (r *MemoryRepository) ListUsers(
	_ context.Context,
	pageSize int32, // pageSize
	_ string, // pageToken
) ([]*domain.User, string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if len(r.users) == 0 {
		return nil, "", domain.NewErrorNotFound("user not found", nil)
	}
	users := make([]*domain.User, 0, len(r.users))
	count := int32(0)
	for _, user := range r.users {
		if user.DeleteTime.IsZero() {
			users = append(users, user)
		}
		count++
		if count >= pageSize {
			break
		}
	}
	return users, "", nil
}

func (r *MemoryRepository) UpdateUser(
	_ context.Context,
	u *domain.User,
) (*domain.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	user, exists := r.users[u.Name]
	if !exists {
		return nil, domain.NewErrorNotFound("user not found", nil)
	}
	userCopy, err := user.Copy()
	if err != nil {
		return nil, domain.NewErrorInternal("failed to copy user", err)
	}
	userCopy.CreateTime = u.CreateTime
	userCopy.UpdateTime = time.Now().UTC()
	r.users[userCopy.Name] = userCopy
	return nil, domain.NewErrorNotFound("user not found", nil)
}

func (r *MemoryRepository) DeleteUser(
	_ context.Context,
	s string, // name
) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	user, exists := r.users[s]
	if !exists || !user.DeleteTime.IsZero() {
		return domain.NewErrorNotFound("user not found", nil)
	}
	user.DeleteTime = time.Now().UTC()
	return nil
}
