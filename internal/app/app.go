package app

import (
	grpcapp "auth/internal/app/grpc"
	"auth/internal/services/auth"
	"auth/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, StoragePath string, TokenTTL time.Duration) *App {
	storage, err := sqlite.New(StoragePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, TokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)
	return &App{grpcApp}
}
