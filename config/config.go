package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	Port              int    `env:"PORT" envDefault:"8080"`
	Logger            Logger `envPrefix:"LOG_"`
	AccuWeatherApiKey string `env:"ACCU_WEATHER_API_KEY"`
	OpenWeatherApiKey string `env:"OPEN_WEATHER_API_KEY"`
	WeatherApiApiKey  string `env:"WEATHER_API_API_KEY"`
	ApiNinjasApiKey   string `env:"API_NINJAS_API_KEY"`
}

type Logger struct {
	Level string `env:"LEVEL" envDefault:"INFO"`
}

func ParseEnv() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		panic("invalid config")
	}

	return &cfg
}
