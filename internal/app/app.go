package app

import (
	"log/slog"
	"time"

	"github.com/Dot-Space/auth_service/config"
	grpc_app "github.com/Dot-Space/auth_service/internal/app/grpc"
	"github.com/Dot-Space/auth_service/internal/auth"
	"github.com/Dot-Space/auth_service/internal/db"
)

type App struct {
	GRPCServer *grpc_app.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	jwtSecret string,
	tokenTTL time.Duration,
	dbConfig config.DBConfig,
) *App {
	db, err := db.New(&dbConfig)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, db, tokenTTL, jwtSecret)

	grpcApp := grpc_app.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
