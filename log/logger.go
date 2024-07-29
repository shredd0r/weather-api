package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"weather-api/config"
)

var std = logrus.New()

func NewLogger(cfg config.Logger) Logger {
	log := logrus.StandardLogger()
	level, err := logrus.ParseLevel(cfg.Level)

	if err != nil {
		panic(fmt.Sprintf("not exist log level: %s", cfg.Level))
	}

	log.SetLevel(level)
	return log
}

func init() {
	std.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

type Logger interface {
	logrus.FieldLogger
}
