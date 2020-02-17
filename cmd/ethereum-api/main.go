package main

import (
	"context"
	"ethereum-api/pkg/app"
	"ethereum-api/pkg/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := app.New()

	ctx, cancel := context.WithCancel(context.Background())
	cancelOnInterrupt(ctx, cancel)

	server.Run(ctx, app)
}

func cancelOnInterrupt(ctx context.Context, cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-ctx.Done():
		case <-c:
			cancel()
		}
	}()
}
