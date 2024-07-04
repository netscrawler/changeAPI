package main

import (
	"context"
	"converter/internal/app"
	"converter/internal/config"
	"converter/internal/lib/logger/handlers/slogpretty"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envlocal = "local"
	envdev   = "dev"
	envprod  = "prod"
)

func main() {

	cfg := config.MustLoad()

	fmt.Println(cfg)

	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.Any("cfg", cfg))
	applicaton := app.New(log, cfg.GRPC.Port)
	go applicaton.GRPCSrv.MustRun()

	// TODO: инициализировать приложение app

	//TODO:redis
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()

	err := client.Set(ctx, "key", "value1", 0).Err()
	if err != nil {
		log.Info("err")
	}

	val, err := client.Get(ctx, "key").Result()
	if err != nil {
		log.Info("err")
	}
	log.Info("key", slog.Any("val", val))

	// TODO: запустить gRPC сервер приложения

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("stoping application", slog.String("signal", sign.String()))
	applicaton.GRPCSrv.Stop()
	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envlocal:
		log = setupPrettySlog()
	case envdev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envprod:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
