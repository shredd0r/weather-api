package config

import (
	"github.com/caarlos0/env/v10"
	"time"
)

type Config struct {
	Port               int                `env:"PORT" envDefault:"8080"`
	Logger             Logger             `envPrefix:"LOG_"`
	Redis              Redis              `envPrefix:"REDIS_"`
	ExpirationDuration ExpirationDuration `envPrefix:"EXPIRATION_"`
	ApiKeys            ApiKeys            `envPrefix:"API_KEY_"`
}

type Redis struct {
	Address string `env:"ADDRESS"`
}

type Logger struct {
	Level string `env:"LEVEL" envDefault:"INFO"`
}

type ExpirationDuration struct {
	WeatherInfo time.Duration `env:"WEATHER_INFO" envDefault:"3h"`
}

type ApiKeys struct {
	AccuWeatherApiKey string `env:"ACCU_WEATHER"`
	OpenWeatherApiKey string `env:"OPEN_WEATHER"`
	WeatherApiApiKey  string `env:"WEATHER_API"`
	ApiNinjasApiKey   string `env:"API_NINJAS"`
}

func ParseEnv() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		panic("invalid config")
	}

	return &cfg
}
