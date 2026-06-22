package handlers

import (
	"auth/internal/service"
	"context"
	"errors"

	"buf.build/go/protovalidate"
	"github.com/dev-yeva/auth-service/protos/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface { // in service package
	Login(ctx context.Context, email, password string, appId int64) (token string, err error)
	Register(ctx context.Context, email, password string) (userId int64, err error)
	IsAdmin(ctx context.Context, userId int64) (isAdmin bool, err error)
}

type AuthHandler struct {
	gen.UnimplementedAuthServer
	auth Auth
}

func RegisterServer(gRPCServer *grpc.Server, auth Auth) {
	gen.RegisterAuthServer(gRPCServer, &AuthHandler{auth: auth})
}

func (a *AuthHandler) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userId, err := a.auth.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &gen.RegisterResponse{UserId: userId}, nil
}

func (a *AuthHandler) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := a.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int64(req.GetAppId()))
	if err != nil {

		if errors.Is(err, service.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		if errors.Is(err, service.ErrInvalidPassword) {
			return nil, status.Error(codes.InvalidArgument, "invalid password")
		}

		if errors.Is(err, service.ErrAppNotFound) {
			return nil, status.Error(codes.NotFound, "app not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &gen.LoginResponse{Token: token}, nil
}

func (a *AuthHandler) IsAdmin(ctx context.Context, req *gen.IsAdminRequest) (*gen.IsAdminResponse, error) {

	isAdmin, err := a.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &gen.IsAdminResponse{IsAdmin: isAdmin}, nil
}
