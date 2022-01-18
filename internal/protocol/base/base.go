package base

import (
	"time"

	"github.com/sxwebdev/protocol-template/internal/config"
	"github.com/sxwebdev/protocol-template/internal/model"
	"github.com/tkcrm/modules/logger"
)

// IBase - base interface for all protocols
type IBase interface {
	ParseData(conn *Conn) error
	Decode(conn *Conn, bs []byte) (model.Locations, error)
	SendCommands(conn *Conn, commands []interface{}) error
}

// Base ...
type Base struct {
	Config       *config.Config
	Logger       logger.Logger
	Manufacturer *model.Manufacturer
	IdleTimeout  time.Duration
	ProtocolType string
}

// New ...
func New(c *config.Config, l logger.Logger, m *model.Manufacturer, protocolType string) *Base {
	loggerExtendedFields := []interface{}{"protocol_type", protocolType}
	return &Base{
		Config:       c,
		Logger:       l.With(loggerExtendedFields...),
		Manufacturer: m,
		ProtocolType: protocolType,
	}
}

func (s *Base) Reply(conn *Conn, data []byte) error {
	_, err := conn.Conn.Write(data)
	return err
}

// AddData ...
func (b *Base) AddData(conn *Conn, data []*model.Location) error {
	if len(data) == 0 {
		return nil
	}

	b.Logger.Debugf("%+v", data[len(data)-1])

	return nil
}

// Используется для дебага протокола
func (s *Base) PrintLastLocTime(locations model.Locations) {
	if s.Config.ENV == "dev" && len(locations) > 0 {
		lastItem := locations[len(locations)-1]
		s.Logger.Debugf("last time: %s", lastItem.Timestamp.String())
	}
}
