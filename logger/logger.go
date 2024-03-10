package logger

import (
	"github.com/sirupsen/logrus"
	"weather-api/config"
)

// TODO
func NewLogger(cfg config.Logger) *logrus.Logger {
	log := logrus.StandardLogger()
	log.SetLevel(logrus.Level(6))

	return log
}
