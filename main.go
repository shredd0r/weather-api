package main

import (
	"context"
	"github.com/shredd0r/weather-api/config"
	"github.com/shredd0r/weather-api/log"
	"github.com/shredd0r/weather-api/server"
)

func main() {
	cfg := config.ParseEnv()
	logger := log.NewLogger(cfg.Logger)

	srv := server.NewServer(cfg, logger)

	srv.Start(context.Background())
}
