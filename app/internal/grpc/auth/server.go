package auth

import (
	"auth/internal/services/auth"
	"auth/protos/gen/go/sso"
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyValue = 0
)

type serverAPI struct {
	sso.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	Register(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

func Register(gRPC *grpc.Server, auth Auth) {
	sso.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *sso.LoginRequest) (*sso.LoginResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.AppId == emptyValue {
		return nil, status.Error(codes.InvalidArgument, "AppId is required")
	}

	///////////////////
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {

		return nil, status.Error(codes.Internal, "Login failed")
	}

	return &sso.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	userID, err := s.auth.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "Register failed")
	}

	return &sso.RegisterResponse{UserId: userID}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *sso.IsAdminRequest) (*sso.IsAdminResponse, error) {
	if req.GetUserId() == emptyValue {
		return nil, status.Error(codes.Internal, "AppId failed")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "Register failed")
	}

	return &sso.IsAdminResponse{IsAdmin: isAdmin}, nil
}
