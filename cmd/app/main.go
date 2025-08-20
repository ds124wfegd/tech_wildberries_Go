// entry point to app :)
package main

import (
	"github.com/ds124wfegd/tech_wildberries_Go/config"
	"github.com/ds124wfegd/tech_wildberries_Go/internal/appServer"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter)) // JSON format for logging

	viperInstance, err := config.LoadConfig() // creating a config entity
	if err != nil {
		logrus.Fatalf("Cannot load config. Error: {%s}", err.Error()) // handling errors related to reading the config
	}

	cfg, err := config.ParseConfig(viperInstance)
	if err != nil {
		logrus.Fatalf("Cannot parse config. Error: {%s}", err.Error()) // handling errors related to parsing config
	}

	appServer.NewServer(cfg) // creating server
}
