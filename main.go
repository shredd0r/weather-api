package main

import (
	"context"
	"weather-api/client/http"
	"weather-api/client/redis"
	"weather-api/config"
	"weather-api/dto"
	"weather-api/log"
	"weather-api/provider"
	"weather-api/service"
	"weather-api/storage"
)

func main() {
	cfg := config.ParseEnv()
	logger := log.NewLogger(cfg.Logger)
	hc := http.NewHttpClient(logger)
	ac := http.NewAccuWeatherClient(logger, hc, cfg.ApiKeys.AccuWeatherApiKey)
	oc := http.NewOpenWeatherClient(logger, hc, cfg.ApiKeys.OpenWeatherApiKey)
	ap := provider.NewOpenWeatherProvider(logger, oc)
	rc := redis.NewClient(&cfg.Redis)
	rws := storage.NewRedisWeatherStorage(rc)
	rls := storage.NewRedisLocationStorage(rc)
	anc := http.NewApiNinjasClient(logger, hc, cfg.ApiKeys.ApiNinjasApiKey)
	ls := service.NewLocationService(logger, rls, ac, anc)
	as := service.NewWeatherService(logger, dto.WeatherForecasterAccuWeather, ls, ap, rws)

	resp, _ := as.DailyWeather(context.Background(), dto.WeatherRequestDto{
		Coords: &dto.Coords{
			Latitude:  50.000691,
			Longitude: 36.215194,
		},
		Locale: "uk-ua",
		Unit:   dto.UnitMetric,
	})

	for _, w := range *resp {
		logger.Info(w.String())
	}

	//port := "8080"

	//srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	//
	//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	//http.Handle("/query", srv)
	//
	//log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	//log.Fatal(http.ListenAndServe(":"+port, nil))
}
