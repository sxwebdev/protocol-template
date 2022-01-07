package protocol

import (
	"os"

	"github.com/sxwebdev/protocol-template/internal/config"
	"github.com/sxwebdev/protocol-template/internal/model"
	"github.com/sxwebdev/protocol-template/internal/protocol/base"
	"github.com/sxwebdev/protocol-template/internal/protocol/newprotocol"
	"github.com/tkcrm/modules/logger"
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

	base_config := &base.Config{
		ProtocolType: "newprotocol",
		Port:         c.ProtocolPort,
		ENV:          c.ENV,
	}
	_base := base.New(base_config, l, &model.Manufacturer{})
	protocol.Protocols[base_config.ProtocolType] = newprotocol.New(_base)

	return protocol, nil
}

// StartServers запускает все сервера для всех протоколов
func (p *Protocol) StartServers(c chan os.Signal) {
	for protocol_name := range p.Protocols {
		if err := p.Protocols[protocol_name].StartServer(); err != nil {
			p.errProtocol(err, c)
		}
	}
}

// StopServers останаливает сервера для всех протоколов
func (p *Protocol) StopServers() {
	for protocol_name := range p.Protocols {
		if err := p.Protocols[protocol_name].StopServer(); err != nil {
			p.logger.Errorf("Stop server %s error: %v", protocol_name, err)
		}
	}
}

// errProtocol отправляет в канал сообщение о том, что произошла ошибка
// и микросервис будет перезапущен
func (p *Protocol) errProtocol(err error, c chan os.Signal) {
	p.logger.Errorf("%+v", err)
	c <- os.Interrupt
}
