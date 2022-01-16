package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sxwebdev/protocol-template/internal/config"
	"github.com/sxwebdev/protocol-template/internal/listenerserver"
	"github.com/sxwebdev/protocol-template/internal/protocol"
	"github.com/sxwebdev/protocol-template/internal/protocol/base"
	"github.com/tkcrm/modules/logger"
	"google.golang.org/grpc"
)

// Server ...
type Server struct {
	Config          *config.Config
	logger          logger.Logger
	grpc            *grpc.Server
	mx              sync.RWMutex
	listenerServers map[string]*listenerserver.ListenerServer
}

// Start server
func Start(l logger.Logger) error {
	s := &Server{
		logger:          l,
		listenerServers: make(map[string]*listenerserver.ListenerServer),
	}

	// Read configuration and envirioments
	config := config.New()
	if err := config.Validate(); err != nil {
		return fmt.Errorf("configuration params validation error: %v", err)
	}
	s.Config = config

	// Init protocols
	pt, err := protocol.New(config, l)
	if err != nil {
		return err
	}
	sigCh := make(chan os.Signal, 1)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	for protocol_name, protocol := range pt.Protocols {
		go func(protocol_name string, protocol base.IBase) {
			port := config.Protocols[protocol_name]
			// TODO реализовать UDP
			ls, err := listenerserver.New(l, protocol_name, port, listenerserver.TypeTCP, protocol)
			if err != nil {
				s.logger.Errorf("start server %s error: %+v", protocol_name, err)
				sigCh <- os.Interrupt
			}
			go func() {
				defer wg.Done()
				wg.Add(1)

				if err := ls.AcceptConnections(ctx); err != nil {
					s.logger.Errorf("accept connection error: %+v", err)
				}
			}()
			s.mx.Lock()
			s.listenerServers[protocol_name] = ls
			s.mx.Unlock()
		}(protocol_name, protocol)
	}

	// Start GRPC server
	go func() {
		if err := s.newGRPC(); err != nil {
			s.logger.Fatal(err)
		}
	}()

	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	cancel()

	// Stop listener servers
	for _, ls := range s.listenerServers {
		if err := ls.StopServer(); err != nil {
			s.logger.Errorf("%+v", err)
		}
	}

	wg.Wait()
	os.Exit(0)

	return nil
}
