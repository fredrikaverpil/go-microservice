package handler

import (
	"errors"

	"github.com/fredrikaverpil/go-microservice/internal/domain"
	pb "github.com/fredrikaverpil/go-microservice/internal/proto/gen/go/gomicroservice/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Domain to Proto conversions
func toProtoUser(user *domain.User) *pb.User {
	return &pb.User{
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		CreateTime:  timestamppb.New(user.CreateTime),
		UpdateTime:  timestamppb.New(user.UpdateTime),
	}
}

// Proto to Domain conversions
func toDomainUser(pbUser *pb.User) *domain.User {
	return &domain.User{
		Name:        pbUser.Name,
		DisplayName: pbUser.DisplayName,
		Email:       pbUser.Email,
		CreateTime:  pbUser.GetCreateTime().AsTime(),
		UpdateTime:  pbUser.GetUpdateTime().AsTime(),
	}
}

// checkTransientError checks for common transient errors that can occur in any operation.
// This should be called before checking specific operation errors.
func checkTransientError(err error) error {
	var customErr *domain.Error
	if !errors.As(err, &customErr) {
		return nil
	}

	switch customErr.Type {
	case domain.Timeout:
		return status.Error(codes.DeadlineExceeded, customErr.Message)
	case domain.Unavailable:
		return status.Error(codes.Unavailable, customErr.Message)
	case domain.ResourceExhausted:
		return status.Error(codes.ResourceExhausted, customErr.Message)
	default:
		return nil
	}
}

// toCreateUserError converts internal errors to gRPC errors following AIP-133.
// Valid error codes for Create methods:
// - InvalidArgument: Client specified invalid/malformed argument.
// - AlreadyExists: The resource already exists.
// - Internal: All other errors are mapped to Internal.
func toCreateUserError(err error) error {
	if transientErr := checkTransientError(err); transientErr != nil {
		return transientErr
	}

	var customErr *domain.Error
	if !errors.As(err, &customErr) {
		return status.Error(codes.Internal, "internal error")
	}

	switch customErr.Type {
	case domain.AlreadyExists:
		return status.Error(codes.AlreadyExists, customErr.Message)
	case domain.InvalidInput:
		return status.Error(codes.InvalidArgument, customErr.Message)
	default:
		return status.Error(codes.Internal, customErr.Message)
	}
}

// toGetUserError converts internal errors to gRPC errors following AIP-131.
// Valid error codes for Get methods:
// - NotFound: The resource was not found.
// - Internal: All other errors are mapped to Internal.
func toGetUserError(err error) error {
	if transientErr := checkTransientError(err); transientErr != nil {
		return transientErr
	}

	var customErr *domain.Error
	if !errors.As(err, &customErr) {
		return status.Error(codes.Internal, "internal error")
	}

	switch customErr.Type {
	case domain.NotFound:
		return status.Error(codes.NotFound, customErr.Message)
	default:
		return status.Error(codes.Internal, customErr.Message)
	}
}

// toListUsersError converts internal errors to gRPC errors following AIP-132.
// Valid error codes for List methods:
// - InvalidArgument: Client specified invalid argument like invalid page token.
// - Internal: All other errors are mapped to Internal.
func toListUsersError(err error) error {
	if transientErr := checkTransientError(err); transientErr != nil {
		return transientErr
	}

	var customErr *domain.Error
	if !errors.As(err, &customErr) {
		return status.Error(codes.Internal, "internal error")
	}

	switch customErr.Type {
	case domain.InvalidInput:
		return status.Error(codes.InvalidArgument, customErr.Message)
	default:
		return status.Error(codes.Internal, customErr.Message)
	}
}

// toUpdateUserError converts internal errors to gRPC errors following AIP-134.
// Valid error codes for Update methods:
// - InvalidArgument: Client specified invalid argument.
// - NotFound: The resource was not found.
// - Internal: All other errors are mapped to Internal.
func toUpdateUserError(err error) error {
	if transientErr := checkTransientError(err); transientErr != nil {
		return transientErr
	}

	var customErr *domain.Error
	if !errors.As(err, &customErr) {
		return status.Error(codes.Internal, "internal error")
	}

	switch customErr.Type {
	case domain.NotFound:
		return status.Error(codes.NotFound, customErr.Message)
	case domain.InvalidInput:
		return status.Error(codes.InvalidArgument, customErr.Message)
	default:
		return status.Error(codes.Internal, customErr.Message)
	}
}

// toDeleteUserError converts internal errors to gRPC errors following AIP-135.
// Valid error codes for Delete methods:
// - NotFound: The resource was not found.
// - Internal: All other errors are mapped to Internal.
func toDeleteUserError(err error) error {
	if transientErr := checkTransientError(err); transientErr != nil {
		return transientErr
	}

	var customErr *domain.Error
	if !errors.As(err, &customErr) {
		return status.Error(codes.Internal, "internal error")
	}

	switch customErr.Type {
	case domain.NotFound:
		return status.Error(codes.NotFound, customErr.Message)
	default:
		return status.Error(codes.Internal, customErr.Message)
	}
}
