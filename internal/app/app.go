package app

import (
	grpcapp "converter/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	tokenTTl time.Duration) *App {
	//TODO: инициализировать хранилище
	//TODO: init convert service
	grpcApp := grpcapp.New(log, grpcPort)
	return &App{
		GRPCSrv: grpcApp,
	}
}
