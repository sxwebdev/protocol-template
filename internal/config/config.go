package config

import (
	"os"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sakirsensoy/genv"
	"github.com/sakirsensoy/genv/dotenv"
	"github.com/tkcrm/modules/utils"
)

// Config ...
type Config struct {
	APPName   string
	APPMSName string
	APPHost   string
	Protocols map[string]uint16
	GRPCPort  string
	ENV       string
}

// New ...
func New() *Config {

	env := utils.GetDefaultString(os.Getenv("ENV"), "dev")

	if env == "dev" {
		dotenv.Load()
	}

	return &Config{
		APPName:   utils.GetDefaultString(genv.Key("APP_NAME").String(), "UNDEFINED_APP_NAME"),
		APPMSName: genv.Key("APP_MS_NAME").String(),
		APPHost:   genv.Key("APP_HOST").String(),
		Protocols: map[string]uint16{
			"newprotocol": 34900,
			// "galileosky": 35100,
			// "teltonika":           35110,
			"adm": 35120,
			// "sat_sol":             35130,
			// "navtelecom": 35140,
			// "arusnavi":            35150,
			// "fort_monitor":        35160,
			// "ruptela":             35170,
			// "egts":                35180,
			// "wialon_ips":          35190,
			// "wialon_retranslator": 35200,
		},
		GRPCPort: utils.GetDefaultString(genv.Key("GRPC_PORT").String(), "9001"),
		ENV:      env,
	}
}

// Validate ...
func (c *Config) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.APPName, validation.Required),
		validation.Field(&c.APPMSName, validation.Required),
		validation.Field(&c.APPHost, validation.Required),
		validation.Field(&c.GRPCPort, validation.Required),
		validation.Field(&c.ENV, validation.Required, validation.In("dev", "stage", "prod")),
	)
}
