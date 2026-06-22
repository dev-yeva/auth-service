package main

import (
	"auth/app"
	"auth/internal/config"
	"auth/lib"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/lmittmann/tint"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	config := config.MustLoad()

	logger := lib.Must(CreateLogger(envLocal))

	application := app.New(logger, config.GRPC.Port, config.StoragePath, config.TokenTTL)
	go application.MustRun()

	handleInterruption(application)
}

func CreateLogger(env string) (*slog.Logger, error) {
	var handler slog.Handler

	switch env {
	case envLocal:
		handler = tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug})
	case envProd:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	default:
		return nil, fmt.Errorf("env must be either '%s' or '%s'", envLocal, envProd)
	}

	return slog.New(handler), nil
}

func handleInterruption(application *app.App) {
	stop_signals := make(chan os.Signal, 1)
	signal.Notify(stop_signals, syscall.SIGINT, syscall.SIGTERM)
	os_signal := <-stop_signals
	application.Stop(slog.Any("os_signal", os_signal))
}
