package service

import (
	pb "auth/genproto/AuthService"
	logger "auth/pkgLogger"
	"auth/storage/postgres"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

type UserService struct {
	pb.UnimplementedAuthUserServiceServer
	User   *postgres.AuthUser
	Logger *slog.Logger
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		User:   postgres.NewAuthUser(db),
		Logger: logger.NewLogger(),
	}
}

func (s *UserService) GetAllUser(ctx context.Context, req *pb.Void) (*pb.GetAllResponse, error) {
	resp, err := s.User.GetAllUser(req)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("error getting data: %v", err))
		return nil, err
	}

	return resp, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateRequest) (*pb.Status, error) {
	resp, err := s.User.UpdateUser(req)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Error updating all user input data: %v", err))
		return nil, err
	}

	return resp, nil
}

func (s *UserService) UpdateUserType(ctx context.Context, req *pb.TypeUserRequest) (*pb.Status, error) {
	resp, err := s.User.UpdateUserType(req)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Error updating user user type: %v", err))
		return nil, err
	}
	return resp, nil
}

func (s *UserService) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	resp, err := s.User.GetUsers(req)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Error getting data with offset and limit: %v", err))
		return nil, err
	}
	return resp, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	resp, err := s.User.DeleteUser(req)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Error deleting user: %v", err))
		return nil, err
	}
	return resp, nil
}

func (s UserService) IdCheck(ctx context.Context, req *pb.Id) (*pb.Response, error) {
	resp, err := s.User.IdCheck(req)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Error checking user id: %v", err))
		return &pb.Response{B: false}, err
	}
	return resp, nil
}


func (s UserService) SearchName(ctx context.Context, req *pb.Id)(*pb.Name,error){
	resp,err:=s.User.SearchName(req)
	if err!=nil{
		return nil,err
	}
	return resp,nil
}