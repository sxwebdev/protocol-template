package base

import (
	"net"
	"sync"
	"time"

	"github.com/sxwebdev/protocol-template/internal/model"
	"github.com/tkcrm/modules/logger"
)

// IBase - base interface for all protocols
type IBase interface {
	StartServer() error
	StopServer() error
	Reply(conn *Conn, data []byte) error
	SendCommands(conn *Conn, commands []interface{}) error
	ConvertData(interface{}) ([]*model.Data, error)
}

// Base ...
type Base struct {
	Config       *Config
	Logger       logger.Logger
	Manufacturer *model.Manufacturer
	IdleTimeout  time.Duration

	mx    sync.RWMutex
	Conns map[*Conn]struct{}
}

// New ...
func New(c *Config, l logger.Logger, m *model.Manufacturer) *Base {
	loggerExtendedFields := []interface{}{"protocol_type", c.ProtocolType}
	return &Base{
		Config:       c,
		Logger:       l.With(loggerExtendedFields...),
		Manufacturer: m,
		Conns:        make(map[*Conn]struct{}),
	}
}

func (b *Base) NewConn(conn net.Conn) *Conn {
	return &Conn{
		Conn:        conn,
		IdleTimeout: time.Minute * 5,
	}
}

// AddConn ...
func (b *Base) AddConn(conn *Conn) {
	b.mx.Lock()
	b.Conns[conn] = struct{}{}
	b.mx.Unlock()
}

// DeleteConn ...
func (b *Base) DeleteConn(conn *Conn) {
	if err := conn.Conn.Close(); err != nil {
		b.Logger.Errorf("Close connection error: %+v", err)
	}
	b.mx.Lock()
	delete(b.Conns, conn)
	b.mx.Unlock()
}

// AddData ...
func (b *Base) AddData(conn *Conn, data []*model.Data) error {
	// TODO ...
	return nil
}
