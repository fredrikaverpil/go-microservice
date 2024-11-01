package handler

import (
	"context"
	"strings"

	"github.com/bufbuild/protovalidate-go"
	"github.com/fredrikaverpil/go-microservice/internal/domain"
	"github.com/fredrikaverpil/go-microservice/internal/port"
	pb "github.com/fredrikaverpil/go-microservice/internal/proto/gen/go/gomicroservice/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCHandler struct {
	pb.UnimplementedUserServiceServer
	userService port.UserService
	validator   *protovalidate.Validator
}

func NewGRPCHandler(userService port.UserService, validator *protovalidate.Validator) *GRPCHandler {
	return &GRPCHandler{
		userService: userService,
		validator:   validator,
	}
}

// CreateUser implements AIP-133
func (h *GRPCHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	// Validate the request
	// TODO: validate fields: clear, validate required using https://github.com/einride/aip-go
	if err := h.validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if req.GetUser() == nil {
		return nil, status.Error(codes.InvalidArgument, "user is required")
	}

	// Convert
	user := &domain.User{
		DisplayName: req.GetUser().DisplayName,
		Email:       req.GetUser().Email,
	}
	// If user_id is provided, use it
	if req.UserId != "" {
		user.Name = "users/" + req.UserId
	}

	// Create
	createdUser, err := h.userService.CreateUser(ctx, user)
	if err != nil {
		return nil, toCreateUserError(err)
	}

	// Convert and return
	return toProtoUser(createdUser), nil
}

// GetUser implements AIP-131
func (h *GRPCHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	// Validate the request
	// TODO: validate fields: clear, validate required using https://github.com/einride/aip-go
	if err := h.validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if !strings.HasPrefix(req.GetName(), "users/") {
		return nil, status.Error(codes.InvalidArgument, "name must start with 'users/'")
	}

	// Get
	user, err := h.userService.GetUser(ctx, req.GetName())
	if err != nil {
		return nil, toGetUserError(err)
	}

	// Convert and return
	return toProtoUser(user), nil
}

// ListUsers implements AIP-132
func (h *GRPCHandler) ListUsers(
	ctx context.Context,
	req *pb.ListUsersRequest,
) (*pb.ListUsersResponse, error) {
	// Validate the request
	if err := h.validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// List
	pageSize := int32(10) // Default page size
	if req.GetPageSize() > 0 {
		pageSize = req.GetPageSize()
	}
	users, nextPageToken, err := h.userService.ListUsers(ctx, pageSize, req.GetPageToken())
	if err != nil {
		return nil, toListUsersError(err)
	}

	// Convert and return
	protoUsers := make([]*pb.User, len(users))
	for i, user := range users {
		protoUsers[i] = toProtoUser(user)
	}
	return &pb.ListUsersResponse{
		Users:         protoUsers,
		NextPageToken: nextPageToken,
	}, nil
}

// UpdateUser implements AIP-134
func (h *GRPCHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	// Validate the request
	// TODO: validate fields: clear, validate required using https://github.com/einride/aip-go
	if err := h.validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if req.GetUser() == nil {
		return nil, status.Error(codes.InvalidArgument, "user is required")
	}
	if !strings.HasPrefix(req.GetUser().GetName(), "users/") {
		return nil, status.Error(codes.InvalidArgument, "user.name must start with 'users/'")
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

// DeleteUser implements AIP-135
func (h *GRPCHandler) DeleteUser(
	ctx context.Context,
	req *pb.DeleteUserRequest,
) (*emptypb.Empty, error) {
	// Validate the request
	// TODO: validate fields: clear, validate required using https://github.com/einride/aip-go
	if err := h.validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if !strings.HasPrefix(req.GetName(), "users/") {
		return nil, status.Error(codes.InvalidArgument, "name must start with 'users/'")
	}

	// Delete
	if err := h.userService.DeleteUser(ctx, req.GetName()); err != nil {
		return nil, toDeleteUserError(err)
	}

	// Return
	return &emptypb.Empty{}, nil
}
