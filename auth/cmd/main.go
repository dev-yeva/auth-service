package main

import (
	"auth/app"
	"auth/internal/config"
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

	logger := MustCreateLogger(envLocal)

	application := app.New(logger, config.GRPC.Port, config.StoragePath, config.TokenTTL)
	go application.MustRun()

	handleInterruption(application)
}

func MustCreateLogger(env string) *slog.Logger {
	var handler slog.Handler

	switch env {
	case envLocal:
		handler = tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug})
	case envProd:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	default:
		panic(fmt.Sprintf("env must be either '%s' or '%s'", envLocal, envProd))
	}
	return slog.New(handler)
}

func handleInterruption(application *app.App) {
	stop_signals := make(chan os.Signal, 1)
	signal.Notify(stop_signals, syscall.SIGINT, syscall.SIGTERM)
	os_signal := <-stop_signals
	application.Stop(slog.Any("os_signal", os_signal))
}
