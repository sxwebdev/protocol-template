package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tkcrm/modules/logger"
	"github.com/sxwebdev/protocol-template/internal/config"
	"github.com/sxwebdev/protocol-template/internal/protocol"
	"google.golang.org/grpc"
)

// Server ...
type Server struct {
	Config   *config.Config
	Logger   logger.Logger
	Protocol *protocol.Protocol
	grpc     *grpc.Server
}

// Start server
func Start(l logger.Logger) error {
	s := &Server{
		Logger: l,
	}

	// Read configuration and envirioments
	config := config.New()
	if err := config.Validate(); err != nil {
		return fmt.Errorf("configuration params validation error: %v", err)
	}
	s.Config = config

	// Init protocols
	c := make(chan os.Signal, 1)
	var err error
	s.Protocol, err = protocol.New(config, l)
	if err != nil {
		return err
	}
	go s.Protocol.StartServers(c)

	// Start GRPC server
	go func() {
		if err := s.newGRPC(); err != nil {
			s.Logger.Fatal(err)
		}
	}()

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	s.Protocol.StopServers()
	os.Exit(1)

	return nil
}
