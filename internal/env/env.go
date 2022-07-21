package env

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gitlab.com/g6834/team41/analytics/internal/cfg"
)

type Environment struct {
	C cfg.Config
}

var E *Environment

var (
	ConfigPath = "CONFIG_PATH"
)

func init() {
	// Get config path from environment variable
	path := os.Getenv(ConfigPath)
	if path == "" {
		path = "config.yaml"
	}

	var err error
	E, err = NewEnvironment(path)
	if err != nil {
		logrus.Panic(fmt.Errorf("failed to load config: %w", err))
	}

	configureLogger()
}

func NewEnvironment(yamlFile string) (*Environment, error) {
	conf, err := cfg.NewConfig(yamlFile)
	if err != nil {
		return nil, err
	}

	return &Environment{C: *conf}, nil
}

func configureLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if E.C.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}
