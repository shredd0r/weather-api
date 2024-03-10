package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	Port        int           `env:"PORT" envDefault:"8080"`
	Logger      Logger        `envPrefix:"LOG_"`
	AccuWeather WeatherApiKey `envPrefix:"ACCU_WEATHER_"`
	OpenWeather WeatherApiKey `envPrefix:"OPEN_WEATHER_"`
	WeatherApi  WeatherApiKey `envPrefix:"WEATHER_API_"`
}

type Logger struct {
	Level string `env:"LEVEL" envDefault:"INFO"`
}

type WeatherApiKey struct {
	ApiKey string `env:"API_KEY"`
}

func ParseEnv() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		panic("invalid config")
	}

	return &cfg
}
