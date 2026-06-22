package app

import (
	"auth/internal/handlers"
	"auth/internal/service"
	"auth/internal/storage/sqlite"
	"auth/lib"
	"fmt"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"
)

type App struct {
	logger     *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(logger *slog.Logger, port int, storagePath string, tokenTTL time.Duration) *App {
	storage := lib.Must(sqlite.New(storagePath))
	authService := service.New(logger, storage, storage, storage, tokenTTL)

	gRPCServer := grpc.NewServer()
	handlers.RegisterServer(gRPCServer, authService)

	return &App{logger: logger, gRPCServer: gRPCServer, port: port}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "app.Run"
	args := []any{
		"op", op,
		"port", a.port,
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.logger.Info("Starting gRPC server", args...)

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop(attrs ...slog.Attr) {
	args := []any{
		"op", "app.Stop",
		"port", a.port,
	}

	for _, attr := range attrs {
		args = append(args, attr)
	}

	a.logger.Info("Stopping gRPC server", args...)
	a.gRPCServer.GracefulStop()
}
