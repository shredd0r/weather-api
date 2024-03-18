package main

import (
	"context"
	"weather-api/client/http"
	"weather-api/client/redis"
	"weather-api/config"
	"weather-api/logger"
	"weather-api/provider"
	"weather-api/repository"
)

func main() {
	ctx := context.Background()
	cfg := config.ParseEnv()
	log := logger.NewLogger(cfg.Logger)
	rClient := redis.NewClient(&ctx, cfg.Redis)

	p := provider.NewAccuWeatherProvider(
		repository.NewRedisCurrentWeatherRepository(log, rClient, cfg.ExpirationDuration),
		http.NewAccuWeatherClient(
			log,
			http.NewHttpClient(log),
			cfg.AccuWeatherApiKey),
		log,
	)

	p.CurrentWeatherInfo("123")

	//port := "8080"
	//
	//srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	//
	//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	//http.Handle("/query", srv)
	//
	//log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	//log.Fatal(http.ListenAndServe(":"+port, nil))
}
