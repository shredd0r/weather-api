package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"weather-api/client"
	"weather-api/dto"
)

func main() {
	logger := NewLogger()

	httpClient := &http.Client{
		Transport: &client.LoggingTransport{
			Logger:    logger,
			Transport: http.DefaultTransport,
		},
	}

	c := client.NewAccuWeatherClient(logger, httpClient)

	c.GetDailyWeatherInfo(
		dto.AccuWeatherRequestDto{
			AccuWeatherBaseRequestDto: dto.AccuWeatherBaseRequestDto{
				AppKey:   "rIvOy0yPABlbkQ9dX1OAwGxwkA4p81hi",
				Language: "uk",
				Details:  true,
				Metric:   true,
			},
			LocationKey: "1212408",
		})

	//port := os.Getenv("PORT")
	//if port == "" {
	//	port = defaultPort
	//}
	//
	//srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	//
	//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	//http.Handle("/query", srv)
	//
	//log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	//log.Fatal(http.ListenAndServe(":"+port, nil))
}

func NewLogger() *logrus.Logger {
	logger := logrus.StandardLogger()
	logger.SetLevel(logrus.Level(6))

	return logger
}
