package main

import (
	"client/cli"
	"client/client"
	"client/config"
	"context"
)

func main() {
	ctx := context.Background()
	cfg := config.Load("config/config.yaml")
	authClient := client.MustNew(
		cfg.Address,
		cfg.Timeout,
		cfg.RetryCount,
	)

	cli.EventLoop(ctx, authClient, cfg.Address)
}
