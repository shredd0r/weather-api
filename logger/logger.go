package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"weather-api/config"
)

func NewLogger(cfg config.Logger) *logrus.Logger {
	log := logrus.StandardLogger()
	level, err := logrus.ParseLevel(cfg.Level)

	if err != nil {
		panic(fmt.Sprintf("not exist log level: %s", cfg.Level))
	}

	log.SetLevel(level)
	return log
}

func StandardLogger() *logrus.Logger {
	return logrus.StandardLogger()
}
