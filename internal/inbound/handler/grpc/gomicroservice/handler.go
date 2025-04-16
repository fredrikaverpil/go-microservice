package gomicroservice

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/fredrikaverpil/go-microservice/internal/core/domain"
	"github.com/fredrikaverpil/go-microservice/internal/core/port"
	gomicroservicev1 "github.com/fredrikaverpil/go-microservice/internal/gen/gomicroservice/v1"
	"go.einride.tech/aip/fieldbehavior"
	"go.einride.tech/aip/resourceid"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCHandler struct {
	gomicroservicev1.UnimplementedUserServiceServer
	userService port.UserService
	validator   protovalidate.Validator
}

func NewGRPCHandler(userService port.UserService, validator protovalidate.Validator) *GRPCHandler {
	return &GRPCHandler{
		userService: userService,
		validator:   validator,
	}
}

// CreateUser implements AIP-133.
func (h *GRPCHandler) CreateUser(
	ctx context.Context,
	req *gomicroservicev1.CreateUserRequest,
) (*gomicroservicev1.User, error) {
	// Validate the request
	if err := h.validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	fieldbehavior.ClearFields(req, annotations.FieldBehavior_OUTPUT_ONLY, annotations.FieldBehavior_IDENTIFIER)
	if err := fieldbehavior.ValidateRequiredFields(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Convert
	user := &domain.User{
		DisplayName: req.GetUser().GetDisplayName(),
		Email:       req.GetUser().GetEmail(),
	}
	if req.GetUserId() != "" {
		// If user_id is provided, use it
		userID := req.GetUserId()
		if err := resourceid.ValidateUserSettable(userID); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		user.Name = "users/" + userID
	} else {
		// Generate a new name using the proper format
		user.Name = "users/" + resourceid.NewSystemGeneratedBase32()
	}

	// Create
	createdUser, err := h.userService.CreateUser(ctx, user)
	if err != nil {
		return nil, toCreateUserError(err)
	}

	// Convert and return
	return toProtoUser(createdUser), nil
}

// GetUser implements AIP-131.
func (h *GRPCHandler) GetUser(
	ctx context.Context,
	req *gomicroservicev1.GetUserRequest,
) (*gomicroservicev1.User, error) {
	// Validate the request
	if err := h.validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := fieldbehavior.ValidateRequiredFields(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	var resourceName gomicroservicev1.UserResourceName
	if err := resourceName.UnmarshalString(req.GetName()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid resource name")
	}
	if resourceName.ContainsWildcard() {
		return nil, status.Error(codes.InvalidArgument, "wildcard not allowed")
	}

	// Get
	user, err := h.userService.GetUser(ctx, req.GetName())
	if err != nil {
		return nil, toGetUserError(err)
	}

	// Convert and return
	return toProtoUser(user), nil
}

// ListUsers implements AIP-132.
func (h *GRPCHandler) ListUsers(
	ctx context.Context,
	req *gomicroservicev1.ListUsersRequest,
) (*gomicroservicev1.ListUsersResponse, error) {
	// Validate the request
	if err := h.validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := fieldbehavior.ValidateRequiredFields(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// List
	pageSize := int32(10) //nolint:mnd // Default page size
	if req.GetPageSize() > 0 {
		pageSize = req.GetPageSize()
	}
	users, nextPageToken, err := h.userService.ListUsers(ctx, pageSize, req.GetPageToken())
	if err != nil {
		return nil, toListUsersError(err)
	}

	// Convert and return
	protoUsers := make([]*gomicroservicev1.User, len(users))
	for i, user := range users {
		protoUsers[i] = toProtoUser(user)
	}
	return &gomicroservicev1.ListUsersResponse{
		Users:         protoUsers,
		NextPageToken: nextPageToken,
	}, nil
}

// UpdateUser implements AIP-134.
func (h *GRPCHandler) UpdateUser(
	ctx context.Context,
	req *gomicroservicev1.UpdateUserRequest,
) (*gomicroservicev1.User, error) {
	// Validate the request
	if err := h.validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	fieldbehavior.ClearFields(req, annotations.FieldBehavior_OUTPUT_ONLY)
	if err := fieldbehavior.ValidateRequiredFields(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	var resourceName gomicroservicev1.UserResourceName
	if err := resourceName.UnmarshalString(req.GetUser().GetName()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid resource name")
	}

	// Convert
	user := toDomainUser(req.GetUser())

	// Update
	updatedUser, err := h.userService.UpdateUser(ctx, user)
	if err != nil {
		return nil, toUpdateUserError(err)
	}

	// Convert and return
	return toProtoUser(updatedUser), nil
}

// DeleteUser implements AIP-135.
func (h *GRPCHandler) DeleteUser(
	ctx context.Context,
	req *gomicroservicev1.DeleteUserRequest,
) (*emptypb.Empty, error) {
	// Validate the request
	if err := h.validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	fieldbehavior.ClearFields(req, annotations.FieldBehavior_OUTPUT_ONLY)
	if err := fieldbehavior.ValidateRequiredFields(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	var resourceName gomicroservicev1.UserResourceName
	if err := resourceName.UnmarshalString(req.GetName()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid resource name")
	}
	if resourceName.ContainsWildcard() {
		return nil, status.Error(codes.InvalidArgument, "wildcard not allowed")
	}

	// Delete
	if err := h.userService.DeleteUser(ctx, req.GetName()); err != nil {
		return nil, toDeleteUserError(err)
	}

	// Return
	return &emptypb.Empty{}, nil
}
