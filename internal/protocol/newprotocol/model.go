package newprotocol

import (
	"github.com/sxwebdev/protocol-template/internal/protocol/base"
)

type NewProtocol struct {
	*base.Base
}

// ProtocolData ...
// edit it
type ProtocolData struct {
	HV   uint8    `tag:"0x01"`
	FW   uint8    `tag:"0x02"`
	Imei [15]byte `tag:"0x03"`
}
