package newprotocol

import (
	"bytes"

	"github.com/sxwebdev/protocol-template/internal/protocol/base"
)

func (s *NewProtocol) Decode(conn *base.Conn, r *bytes.Reader) (interface{}, map[string]interface{}, error) {

	packets := make([]ProtocolData, 0)

	return packets, nil, nil
}
