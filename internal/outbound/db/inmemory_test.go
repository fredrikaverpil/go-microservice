package db

import (
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/fredrikaverpil/go-microservice/internal/core/domain"
	gomicroservicev1 "github.com/fredrikaverpil/go-microservice/internal/inbound/handler/grpc/gen/go/gomicroservice/v1"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gotest.tools/v3/assert"
)

// ignoredTimeFields defines which User fields to ignore in comparisons.
var ignoredTimeFields = cmpopts.IgnoreFields( //nolint:gochecknoglobals // this is just a test file
	domain.User{},
	"CreateTime",
	"UpdateTime",
)

func setupTestRepo(_ *testing.T) *MemoryRepository {
	logger := slog.Default()
	repo := NewMemoryRepository(logger).(*MemoryRepository)
	return repo
}

// isRecentTime returns true if the given time is within the last second.
func isRecentTime(t time.Time) bool {
	now := time.Now().UTC()
	difference := now.Sub(t)
	return difference >= 0 && difference < time.Second
}

// assertValidTimestamps verifies that CreateTime and UpdateTime are:
// - Not zero.
// - In UTC.
// - Recent (within the last second).
// - UpdateTime is not before CreateTime.
func assertValidTimestamps(t *testing.T, user *domain.User) {
	t.Helper()

	assert.Assert(t, !user.CreateTime.IsZero(), "CreateTime should not be zero")
	assert.Assert(t, !user.UpdateTime.IsZero(), "UpdateTime should not be zero")

	assert.Equal(t, user.CreateTime.Location(), time.UTC, "CreateTime should be in UTC")
	assert.Equal(t, user.UpdateTime.Location(), time.UTC, "UpdateTime should be in UTC")

	assert.Assert(t, isRecentTime(user.CreateTime), "CreateTime should be recent")
	assert.Assert(t, isRecentTime(user.UpdateTime), "UpdateTime should be recent")

	assert.Assert(t, !user.UpdateTime.Before(user.CreateTime),
		"UpdateTime should not be before CreateTime")
}

// TestCreateUser tests the Create method following AIP-133 (Create Resource).
func TestCreateUser(t *testing.T) {
	t.Parallel()

	t.Run("success with auto-generated ID", func(t *testing.T) {
		t.Parallel()
		repo := setupTestRepo(t)
		ctx := context.Background()

		inputUser := &domain.User{
			DisplayName: "Test User",
			Email:       "test@example.com",
		}

		created, err := repo.CreateUser(ctx, inputUser)
		assert.NilError(t, err)
		assert.Assert(t, created != nil)
		assert.Assert(t, strings.HasPrefix(created.Name, "users/"))

		// Verify timestamps before comparing the rest of the fields
		assertValidTimestamps(t, created)

		// Create expected user after we know the generated ID
		expectedUser := &domain.User{
			Name:        created.Name, // Use the generated name
			DisplayName: "Test User",
			Email:       "test@example.com",
		}

		assert.DeepEqual(t, expectedUser, created, ignoredTimeFields)
	})

	t.Run("success with provided ID", func(t *testing.T) {
		t.Parallel()
		repo := setupTestRepo(t)
		ctx := context.Background()

		inputUser := &domain.User{
			Name:        "users/test123",
			DisplayName: "Test User",
			Email:       "test@example.com",
		}

		created, err := repo.CreateUser(ctx, inputUser)
		assert.NilError(t, err)

		// Verify timestamps before comparing the rest of the fields
		assertValidTimestamps(t, created)

		expectedUser := &domain.User{
			Name:        "users/test123",
			DisplayName: "Test User",
			Email:       "test@example.com",
		}

		assert.DeepEqual(t, expectedUser, created, ignoredTimeFields)
	})

	t.Run("failure - already exists", func(t *testing.T) {
		t.Parallel()
		repo := setupTestRepo(t)
		ctx := context.Background()

		user := &domain.User{
			Name:        "users/duplicate",
			DisplayName: "Test User",
			Email:       "test@example.com",
		}

		// Create first time
		_, err := repo.CreateUser(ctx, user)
		assert.NilError(t, err)

		// Try to create again
		_, err = repo.CreateUser(ctx, user)
		assert.DeepEqual(t, err, &domain.Error{
			Type:    domain.AlreadyExists,
			Message: "user already exists: users/duplicate",
		})
	})

	t.Run("failure - invalid name format", func(t *testing.T) {
		t.Parallel()
		repo := setupTestRepo(t)
		ctx := context.Background()

		user := &domain.User{
			Name:        "invalid-format", // this should fail as it doesn't follow users/{user} format
			DisplayName: "Test User",
			Email:       "test@example.com",
		}

		_, err := repo.CreateUser(ctx, user)
		assert.Assert(t, err != nil)
		var resourceName gomicroservicev1.UserResourceName
		validationErr := resourceName.UnmarshalString(user.Name)
		assert.Assert(t, validationErr != nil)
	})

	t.Run("failure - missing required fields", func(t *testing.T) {
		t.Parallel()
		repo := setupTestRepo(t)
		ctx := context.Background()

		user := &domain.User{
			Name: "users/test123",
		}

		_, err := repo.CreateUser(ctx, user)
		assert.DeepEqual(t, err, &domain.Error{
			Type:    domain.InvalidInput,
			Message: "display_name is required",
		})
	})
}

// TestGetUser tests the Get method following AIP-131 (Get Resource).
func TestGetUser(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		repo := setupTestRepo(t)
		ctx := context.Background()

		inputUser := &domain.User{
			Name:        "users/test123",
			DisplayName: "Test User",
			Email:       "test@example.com",
		}

		created, err := repo.CreateUser(ctx, inputUser)
		assert.NilError(t, err)

		// Get the user
		retrieved, err := repo.GetUser(ctx, created.Name)
		assert.NilError(t, err)

		// Verify timestamps are preserved
		assert.DeepEqual(t, created.CreateTime, retrieved.CreateTime)
		assert.DeepEqual(t, created.UpdateTime, retrieved.UpdateTime)

		expectedUser := &domain.User{
			Name:        "users/test123",
			DisplayName: "Test User",
			Email:       "test@example.com",
		}

		// We expect exactly the same object that was created
		assert.DeepEqual(t, expectedUser, retrieved, ignoredTimeFields)
	})

	t.Run("failure - not found", func(t *testing.T) {
		t.Parallel()
		repo := setupTestRepo(t)
		ctx := context.Background()

		_, err := repo.GetUser(ctx, "users/nonexistent")
		assert.DeepEqual(t, err, &domain.Error{
			Type:    domain.NotFound,
			Message: "user not found",
		})
	})
}
