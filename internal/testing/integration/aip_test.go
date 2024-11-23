package integration

import (
	"context"
	"encoding/base32"
	"regexp"
	"testing"

	gomicroservicev1 "github.com/fredrikaverpil/go-microservice/internal/gen/gomicroservice/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type aipTests struct{}

var _ gomicroservicev1.ServiceConfigProviders = &aipTests{}

func TestUserService(t *testing.T) {
	gomicroservicev1.TestServices(t, &aipTests{})
}

// NewSystemGenerated returns a new system-generated resource ID encoded as base32 lowercase.
func NewSystemGeneratedBase32() string {
	base32Encoding := base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567").WithPadding(base32.NoPadding)
	regexp := regexp.MustCompile(`^[a-z]([a-z0-9-]{0,61}[a-z0-9])?$`)
	// Retry creating id until it matches the regexp
	for {
		id := uuid.New()
		encodedID := base32Encoding.EncodeToString(id[:])
		if regexp.MatchString(encodedID) {
			return encodedID
		}
	}
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
		// NOTE: the generator provided by resourceid does not follow the regex from
		// https://google.aip.dev/122#resource-id-segments
		//
		// IDGenerator: resourceid.NewSystemGeneratedBase32,
		IDGenerator: NewSystemGeneratedBase32,
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
