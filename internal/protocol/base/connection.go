package base

import (
	"net"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Conn ...
type Conn struct {
	Conn        net.Conn
	IdleTimeout time.Duration
	// Device IMEI
	IMEI string
	// Device SerialNumber. For example fot sattelite solutions
	SerialNumber string
	DeviceID     string
	Hardware     string
	Firmware     string
	// Сюда можно складывать разные параметры для конкретного устройства
	// которые доступны например только в первом пакете
	Params map[string]interface{}
}

// Validate ...
func (c *Conn) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.IMEI, validation.Required),
	)
}
