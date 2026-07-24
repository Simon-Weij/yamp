package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Simon-Weij/yamp/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		slog.Error("failed to listen", "err", err)
		os.Exit(1)
	}

	if err := app.Run(ctx, ln); err != nil {
		slog.Error("application error", "err", err)
		os.Exit(1)
	}
}
