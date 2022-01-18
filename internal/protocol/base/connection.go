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

func (c *Conn) SetIMEI(v string) error {
	c.IMEI = v

	return nil
}

func (c *Conn) SetSerialNumber(v string) {
	c.SerialNumber = v
}

func (c *Conn) SetDeviceId(v string) {
	c.DeviceID = v
}

func (c *Conn) SetHardware(v string) {
	c.Hardware = v
}

func (c *Conn) SetFirmware(v string) {
	c.Firmware = v
}

func (c *Conn) SetParam(key string, value interface{}) {
	if c.Params == nil {
		c.Params = map[string]interface{}{}
	}
	c.Params[key] = value
}

// Validate ...
func (c *Conn) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.IMEI, validation.Required),
	)
}
