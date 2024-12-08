package integration

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	gomicroservicev1 "github.com/fredrikaverpil/go-microservice/internal/gen/gomicroservice/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type aipTests struct{}

var (
	_ gomicroservicev1.ServiceConfigProviders = &aipTests{}

	// counter to keep track of unique IDs.
	idCounter uint64 //nolint:gochecknoglobals // used for generating unique IDs.
)

func TestUserService(t *testing.T) {
	gomicroservicev1.TestServices(t, &aipTests{})
}

func (a aipTests) UserServiceUser(_ *testing.T) *gomicroservicev1.UserServiceUserTestSuiteConfig {
	fx := newFixture()

	return &gomicroservicev1.UserServiceUserTestSuiteConfig{
		Service: func() gomicroservicev1.UserServiceServer {
			return fx.userHandler
		},
		Context: func() context.Context {
			return fx.ctx
		},
		Create: func() *gomicroservicev1.User {
			now := timestamppb.Now()

			return &gomicroservicev1.User{
				Name:        "users/johndoe",
				DisplayName: "John Doe",
				Email:       "johndoe@gmail.com",
				CreateTime:  now,
				UpdateTime:  now,
			}
		},
		IDGenerator: func() string {
			id := atomic.AddUint64(&idCounter, 1)
			return fmt.Sprintf("valid-id-%d", id)
		},
		Update: func() *gomicroservicev1.User {
			now := timestamppb.Now()
			return &gomicroservicev1.User{
				Name:        "users/janedoe",
				DisplayName: "Jane Doe",
				Email:       "janedoe@gmail.com",
				UpdateTime:  now,
			}
		},
		Skip: []string{
			"Update/persisted",
			"Update/update_time",
			"Update/preserve_create_time",
			"Update/invalid_update_mask",
			"List/negative_page_size",
			"List/invalid_page_token",
			"List/negative_pages_size",
		},
	}
}
