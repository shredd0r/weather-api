package main

import (
	"context"
	"weather-api/client/http"
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
	rws := &storage.RedisWeatherStorage{}
	rls := &storage.RedisLocationStorage{}
	anc := http.NewApiNinjasClient(logger, hc, cfg.ApiKeys.ApiNinjasApiKey)
	ls := service.NewLocationService(logger, rls, ac, anc)
	as := service.NewWeatherService(logger, dto.WeatherForecasterAccuWeather, ls, ap, rws)

	resp, err := as.CurrentWeather(context.Background(), dto.WeatherRequestDto{
		Coords: &dto.Coords{
			Latitude:  50.000691,
			Longitude: 36.215194,
		},
		Locale: "uk-ua",
		Unit:   dto.UnitMetric,
	})

	logger.Error(err)
	logger.Info(*resp)

	//port := "8080"

	//srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	//
	//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	//http.Handle("/query", srv)
	//
	//log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	//log.Fatal(http.ListenAndServe(":"+port, nil))
}
