package main

import (
	"os"

	"github.com/sakirsensoy/genv"
	"github.com/sakirsensoy/genv/dotenv"
	"github.com/tkcrm/modules/logger"
	"github.com/tkcrm/modules/utils"
	"github.com/sxwebdev/protocol-template/internal/server"
)

func main() {

	if utils.GetDefaultString(os.Getenv("ENV"), "dev") == "dev" {
		dotenv.Load()
	}

	l := logger.DefaultLogger(
		utils.GetDefaultString(genv.Key("LOG_LEVEL").String(), "info"),
		utils.GetDefaultString(genv.Key("APP_MS_NAME").String(), "undefined_app_ms_name"),
	)
	if err := server.Start(l); err != nil {
		l.Errorf("INIT APP ERROR: %v", err)
	}
}
