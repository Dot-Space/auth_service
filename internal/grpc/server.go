package grpc_server

import (
	"context"

	auth_gen "github.com/Dot-Space/auth_service/gen-proto/gen/go/auth"

	"google.golang.org/grpc"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
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
	) (token string, err error)

	RegisterUser(
		ctx context.Context,
		email string,
		password string,
	) (uid int64, err error)

	CheckToken(
		ctx context.Context,
		token string,
	) (status bool, reason string, err error)

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

}

func (s *serverAPI) Register(
	ctx context.Context,
	in *auth_gen.RegisterRequest,
) (*auth_gen.RegisterResponse, error) {

}

func (s *serverAPI) CheckToken(
	ctx context.Context,
	in *auth_gen.CheckRequest,
) (*auth_gen.CheckResponse, error) {

}

func (s *serverAPI) RefreshToken(
	ctx context.Context,
	in *auth_gen.RefreshRequest,
) (*auth_gen.RefreshResponse, error) {

}
