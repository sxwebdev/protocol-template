package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/sxwebdev/protocol-template/internal/config"
	"github.com/sxwebdev/protocol-template/internal/server"
	"github.com/tkcrm/modules/pkg/cfg"
	"github.com/tkcrm/modules/pkg/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGSTOP, syscall.SIGTERM)
	defer stop()

	// Load configuration
	var config config.Config
	if err := cfg.LoadConfig(&config); err != nil {
		logger.New().Fatalf("could not load configuration: %v", err)
	}

	// Init logger
	l := logger.New(
		logger.WithAppName(config.ServiceName),
	)
	defer l.Sync()

	// Init Server
	srv := server.New(l, &config)

	// Start server
	if err := srv.Start(ctx); err != nil {
		l.Fatalf("start server error: %v", err)
	}
}
