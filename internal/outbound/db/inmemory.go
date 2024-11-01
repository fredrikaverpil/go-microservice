package db

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/fredrikaverpil/go-microservice/internal/core/domain"
	"github.com/fredrikaverpil/go-microservice/internal/core/port"
	"github.com/google/uuid"
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

	// Extract user_id from name or generate new one
	var userID string
	if user.Name != "" {
		parts := strings.Split(user.Name, "/")
		if len(parts) != 2 || parts[0] != "users" {
			return nil, domain.NewErrorInvalidInput("invalid name format", nil)
		}
		userID = parts[1]
	} else {
		userID = uuid.New().String()
		user.Name = fmt.Sprintf("users/%s", userID)
	}

	// Check if user already exists
	if _, exists := r.users[user.Name]; exists {
		return nil, domain.NewErrorAlreadyExists(
			fmt.Sprintf("user already exists: %s", user.Name),
			nil,
		)
	}

	// Validate required fields
	if user.DisplayName == "" {
		return nil, domain.NewErrorInvalidInput("display_name is required", nil)
	}
	if user.Email == "" {
		return nil, domain.NewErrorInvalidInput("email is required", nil)
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
	return r.copyUser(newUser), nil
}

func (r *MemoryRepository) copyUser(user *domain.User) *domain.User {
	if user == nil {
		return nil
	}
	return &domain.User{
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		CreateTime:  user.CreateTime,
		UpdateTime:  user.UpdateTime,
	}
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
