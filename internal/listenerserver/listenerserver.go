package listenerserver

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/sxwebdev/protocol-template/internal/protocol/base"
	"github.com/tkcrm/modules/logger"
)

type ListenerServerType int

const (
	TypeTCP ListenerServerType = iota
	TypeUDP
)

type ListenerServer struct {
	logger       logger.Logger
	protocolType string
	listener     net.Listener
	mx           sync.RWMutex
	protocol     base.IBase
	Conns        map[*base.Conn]struct{}
}

func New(l logger.Logger, protocolType string, port uint16, serverType ListenerServerType, protocol base.IBase) (*ListenerServer, error) {

	loggerExtendedFields := []interface{}{"protocol_type", protocolType}
	l = l.With(loggerExtendedFields...)

	var listenerNetwork string
	switch serverType {
	case TypeTCP:
		listenerNetwork = "tcp"
	case TypeUDP:
		listenerNetwork = "udp"
	default:
		return nil, fmt.Errorf("server type %v is not supported", serverType)
	}

	listener, err := net.Listen(listenerNetwork, fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	l.Infof("server %s started", protocolType)

	return &ListenerServer{
		logger:       l,
		protocolType: protocolType,
		listener:     listener,
		protocol:     protocol,
		Conns:        make(map[*base.Conn]struct{}),
	}, nil
}

// StopServer ...
func (s *ListenerServer) StopServer() error {

	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}

func (s *ListenerServer) NewConn(conn net.Conn) *base.Conn {
	return &base.Conn{
		Conn:        conn,
		IdleTimeout: time.Minute * 5,
		Params:      make(map[string]interface{}),
	}
}

func (s *ListenerServer) AddConn(conn *base.Conn) {
	s.mx.Lock()
	s.Conns[conn] = struct{}{}
	s.mx.Unlock()
}

func (s *ListenerServer) DeleteConn(conn *base.Conn) {
	if err := conn.Conn.Close(); err != nil {
		s.logger.Errorf("Close connection error: %+v", err)
	}
	s.mx.Lock()
	delete(s.Conns, conn)
	s.mx.Unlock()
}

func (s *ListenerServer) AcceptConnections(ctx context.Context) error {
	for {
		connection, err := s.listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				s.logger.Errorf("Error accepting connection %v", err)
				continue
			}
		}

		s.logger.Debugf("Accepted connection from %v", connection.RemoteAddr())

		go func(connection net.Conn) {
			if err := s.handleConnection(ctx, s.NewConn(connection)); err != nil {
				s.logger.Errorf("handleConnection error: %v", err)
			}
		}(connection)
	}
}

func (s *ListenerServer) handleConnection(ctx context.Context, conn *base.Conn) error {

	defer s.DeleteConn(conn)
	s.AddConn(conn)

	done := make(chan error, 1)
	timer := time.NewTimer(conn.IdleTimeout)
	defer timer.Stop()

	for {

		go func(conn *base.Conn, done chan error) {
			err := s.protocol.ParseData(conn)
			done <- err
		}(conn, done)

		select {
		case <-timer.C:
			return fmt.Errorf("stop by timeout")
		case <-ctx.Done():
			return nil
		case err := <-done:
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(conn.IdleTimeout)
			if err != nil {
				return err
			}
		}
	}
}
