package newprotocol

import (
	"github.com/sxwebdev/protocol-template/internal/protocol/base"
)

func (s *NewProtocol) Decode(conn *base.Conn, bs []byte) (interface{}, error) {

	packets := make([]ProtocolData, 0)

	return packets, nil
}
