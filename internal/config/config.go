package config

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Config ...
type Config struct {
	ServiceName string            `json:"SERVICE_NAME" default:"protocol-template"`
	StandName   string            `json:"STAND_NAME" default:"local"`
	Protocols   map[string]uint16 `json:"PROTOCOLS"`
	GrpcAddr    string            `json:"GRPC_ADDR" default:":9001"`
}

// Validate ...
func (c *Config) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.ServiceName, validation.Required),
		validation.Field(&c.StandName, validation.Required, validation.In("local", "dev", "stage", "prod")),
		validation.Field(&c.GrpcAddr, validation.Required),
	)
}
