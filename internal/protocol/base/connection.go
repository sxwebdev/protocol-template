package base

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Conn ...
type Conn struct {
	Conn        net.Conn
	Reader      *bufio.Reader
	IdleTimeout time.Duration
	// Device IMEI
	IMEI string
	// Device SerialNumber. For example fot sattelite solutions
	SerialNumber string
	DeviceID     string
	Hardware     string
	Firmware     string
	Params       map[string]interface{}
}

func NewConn(conn net.Conn) *Conn {
	return &Conn{
		Conn:        conn,
		Reader:      bufio.NewReader(conn),
		IdleTimeout: time.Minute * 5,
		Params:      make(map[string]interface{}),
	}
}

func (c *Conn) CheckIMEI(imei string) error {
	r := regexp.MustCompile(`^\d+$`)
	if !r.MatchString(imei) {
		return fmt.Errorf("invalid imei: %s", imei)
	}

	if len(imei) < 15 || len(imei) > 16 {
		return fmt.Errorf("invalid imei length: %s / len: %d", imei, len(imei))
	}

	return nil
}

func (c *Conn) SetIMEI(v string) error {
	if err := c.CheckIMEI(v); err != nil {
		return err
	}
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

func (c *Conn) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.IMEI, validation.Required),
	)
}
