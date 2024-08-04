package config

import (
	"github.com/caarlos0/env/v10"
	"time"
)

type Config struct {
	Server             Server             `envPrefis:"SERVER_"`
	Logger             Logger             `envPrefix:"LOG_"`
	Redis              Redis              `envPrefix:"REDIS_"`
	ExpirationDuration ExpirationDuration `envPrefix:"EXPIRATION_"`
	ApiKeys            ApiKeys            `envPrefix:"API_KEY_"`
}

type Server struct {
	Port             int  `env:"PORT" envDefault:"8080"`
	PlaygroundEnable bool `env:"PLAYGROUND_ENABLE" envDefault:"true"`
}

type Redis struct {
	Address string `env:"ADDRESS"`
}

type Logger struct {
	Level string `env:"LEVEL" envDefault:"INFO"`
}

type ExpirationDuration struct {
	WeatherInfo time.Duration `env:"WEATHER_INFO" envDefault:"3h"`
	Coords      time.Duration `env:"COORDS" envDefault:"3h"`
	Period      time.Duration `env:"PERIOD" envDefault:"5m"`
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
