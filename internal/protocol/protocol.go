package protocol

import (
	"fmt"

	"github.com/sxwebdev/protocol-template/internal/config"
	"github.com/sxwebdev/protocol-template/internal/model"
	"github.com/sxwebdev/protocol-template/internal/protocol/adm"
	"github.com/sxwebdev/protocol-template/internal/protocol/base"
	"github.com/sxwebdev/protocol-template/internal/protocol/newprotocol"
	"github.com/tkcrm/modules/pkg/logger"
)

// Protocol ...
type Protocol struct {
	logger    logger.Logger
	Protocols map[string]base.IBase
}

// New ...
func New(c *config.Config, l logger.Logger) (*Protocol, error) {
	protocol := &Protocol{
		logger:    l,
		Protocols: make(map[string]base.IBase),
	}

	for protocol_name := range c.Protocols {
		_base := base.New(c, l, &model.Manufacturer{}, protocol_name)
		switch protocol_name {
		case "newprotocol":
			protocol.Protocols[protocol_name] = newprotocol.New(_base)
		case "adm":
			protocol.Protocols[protocol_name] = adm.New(_base)
		default:
			return nil, fmt.Errorf("protocol %s is not supported", protocol_name)
		}
	}

	return protocol, nil
}
