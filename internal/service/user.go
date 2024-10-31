package service

import (
	"context"

	pb "github.com/fredrikaverpil/go-microservice/internal/proto/gen/go/gomicroservice/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	// TODO: Implement creation logic
	return &pb.User{}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	// TODO: Implement get logic
	return &pb.User{}, nil
}

func (s *UserService) ListUsers(
	ctx context.Context,
	req *pb.ListUsersRequest,
) (*pb.ListUsersResponse, error) {
	// TODO: Implement list logic
	return &pb.ListUsersResponse{}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	// TODO: Implement update logic
	return &pb.User{}, nil
}

func (s *UserService) DeleteUser(
	ctx context.Context,
	req *pb.DeleteUserRequest,
) (*emptypb.Empty, error) {
	// TODO: Implement delete logic
	return &emptypb.Empty{}, nil
}
