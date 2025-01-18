package grpc_server

import (
	"context"
	"errors"

	auth_gen "github.com/Dot-Space/auth_service/gen-proto/gen/go/auth"
	"github.com/Dot-Space/auth_service/internal/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Функционал API
type serverAPI struct {
	auth_gen.UnimplementedAuthServer
	auth Auth
}

// Интерфейс сервиса
type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
	) (token string, refreshToken string, err error)

	RegisterUser(
		ctx context.Context,
		email string,
		password string,
	) (uid int64, err error)

	CheckToken(
		ctx context.Context,
		token string,
		tokenType string,
	) (status bool, err error)

	RefreshToken(
		ctx context.Context,
		refresh_token string,
	) (token string, new_refresh_token string, err error)
}

// Функция для привязки (регистрации) сервиса к GRPC серверу с использованием вышеописанного интерфейса
func Register(gRPCServer *grpc.Server, auth Auth) {
	auth_gen.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

// Методы API вместе с валидацией входных данных

func (s *serverAPI) Login(
	ctx context.Context,
	in *auth_gen.LoginRequest,
) (*auth_gen.LoginResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "Email is required!")
	}
	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "Email is required!")
	}

	token, refreshToken, err := s.auth.Login(ctx, in.Email, in.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "Invalid email or password!")
		}

		return nil, status.Error(codes.Internal, "Failed to log in")
	}

	return &auth_gen.LoginResponse{Token: token, RefreshToken: refreshToken}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	in *auth_gen.RegisterRequest,
) (*auth_gen.RegisterResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "Email is required!")
	}
	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "Email is required!")
	}

	uid, err := s.auth.RegisterUser(ctx, in.Email, in.Password)
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "Invalid email or password!")
		}

		return nil, status.Error(codes.Internal, "Failed to log in")
	}

	return &auth_gen.RegisterResponse{Uid: uid}, nil
}

func (s *serverAPI) CheckToken(
	ctx context.Context,
	in *auth_gen.CheckRequest,
) (*auth_gen.CheckResponse, error) {
	if in.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "Token is required!")
	}
	if in.TokenType == "" {
		return nil, status.Error(codes.InvalidArgument, "Token type is required!")
	}

	isValid, err := s.auth.CheckToken(ctx, in.Token, in.TokenType)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid token")
	}

	return &auth_gen.CheckResponse{Status: isValid}, nil
}

func (s *serverAPI) RefreshToken(
	ctx context.Context,
	in *auth_gen.RefreshRequest,
) (*auth_gen.RefreshResponse, error) {
	if in.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "Refresh token is required!")
	}

	newAccessToken, newRefreshToken, err := s.auth.RefreshToken(ctx, in.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Refresh token is invalid!")
	}

	return &auth_gen.RefreshResponse{
		Token:           newAccessToken,
		NewRefreshToken: newRefreshToken,
	}, nil
}
