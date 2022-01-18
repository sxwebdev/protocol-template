package newprotocol

import (
	"errors"

	"github.com/sxwebdev/protocol-template/internal/protocol/base"
)

// New ...
func New(b *base.Base) base.IBase {
	return &NewProtocol{Base: b}
}

func (s *NewProtocol) ParseData(conn *base.Conn) error {

	return errors.New("not implemented")
}

// SendCommand ...
// edit it
func (s *NewProtocol) SendCommands(conn *base.Conn, commands []interface{}) error {
	return errors.New("not implemented")
}
