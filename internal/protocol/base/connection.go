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
}

// Validate ...
func (c *Conn) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.IMEI, validation.Required),
	)
}
