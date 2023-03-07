package server

import (
	"context"
	"sync"

	"github.com/sxwebdev/protocol-template/internal/config"
	"github.com/sxwebdev/protocol-template/internal/listenerserver"
	"github.com/sxwebdev/protocol-template/internal/protocol"
	"github.com/sxwebdev/protocol-template/internal/protocol/base"
	"github.com/tkcrm/modules/pkg/logger"
	"google.golang.org/grpc"
)

// Server ...
type Server struct {
	config          *config.Config
	logger          logger.Logger
	grpc            *grpc.Server
	mx              sync.RWMutex
	listenerServers map[string]*listenerserver.ListenerServer
}

func New(l logger.Logger, cfg *config.Config) *Server {
	s := &Server{
		logger:          l,
		config:          cfg,
		listenerServers: make(map[string]*listenerserver.ListenerServer),
	}

	s.config.Protocols = map[string]uint16{
		"newprotocol": 34900,
		"adm":         35120,
	}

	return s
}

// Start server
func (s *Server) Start(ctx context.Context) error {
	// init protocols
	pt, err := protocol.New(s.config, s.logger)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for protocol_name, protocol := range pt.Protocols {
		go func(protocol_name string, protocol base.IBase) {
			port := s.config.Protocols[protocol_name]
			// TODO реализовать UDP
			ls, err := listenerserver.New(s.logger, protocol_name, port, listenerserver.TypeTCP, protocol)
			if err != nil {
				s.logger.Errorf("start server %s error: %+v", protocol_name, err)
				cancel()
				return
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

	<-ctx.Done()

	// Stop listener servers
	for _, ls := range s.listenerServers {
		if err := ls.StopServer(); err != nil {
			s.logger.Errorf("%+v", err)
		}
	}

	wg.Wait()

	return nil
}
