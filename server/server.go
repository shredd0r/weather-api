package server

import (
	"context"
	"fmt"
	http2 "net/http"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/shredd0r/weather-api/client/http"
	"github.com/shredd0r/weather-api/client/redis"
	"github.com/shredd0r/weather-api/config"
	"github.com/shredd0r/weather-api/dto"
	"github.com/shredd0r/weather-api/graph"
	"github.com/shredd0r/weather-api/graph/api"
	"github.com/shredd0r/weather-api/log"
	"github.com/shredd0r/weather-api/provider"
	"github.com/shredd0r/weather-api/service"
	"github.com/shredd0r/weather-api/storage"
	"github.com/shredd0r/weather-api/task"
)

type Server struct {
	cfg              *config.Config
	wg               *sync.WaitGroup
	logger           log.Logger
	graphqlApi       *api.WeatherGraphqlApi
	listCleanerTasks []task.Task
}

func NewServer(cfg *config.Config, logger log.Logger) *Server {
	httpClient := http.NewHttpClient(logger)
	apiNinjasClient := http.NewApiNinjasClient(logger, httpClient, cfg.ApiKeys.ApiNinjasApiKey)
	accuWeatherClient := http.NewAccuWeatherClient(logger, httpClient, cfg.ApiKeys.AccuWeatherApiKey)
	openWeatherClient := http.NewOpenWeatherClient(logger, httpClient, cfg.ApiKeys.OpenWeatherApiKey)

	accuWeatherProvider := provider.NewAccuWeatherProvider(logger, accuWeatherClient)
	openWeatherProvider := provider.NewOpenWeatherProvider(logger, openWeatherClient)

	redisClient := redis.NewClient(&cfg.Redis)
	redisWeatherStorage := storage.NewRedisWeatherStorage(logger, redisClient)
	redisLocationStorage := storage.NewRedisLocationStorage(logger, redisClient)

	locationService := service.NewLocationService(logger, redisLocationStorage, accuWeatherClient, apiNinjasClient)
	accuWeatherService := service.NewWeatherService(logger, dto.WeatherForecasterAccuWeather, locationService, accuWeatherProvider, redisWeatherStorage)
	openWeatherService := service.NewWeatherService(logger, dto.WeatherForecasterOpenWeather, locationService, openWeatherProvider, redisWeatherStorage)

	return &Server{
		cfg:              cfg,
		logger:           logger,
		wg:               &sync.WaitGroup{},
		graphqlApi:       api.NewWeatherGraphqlApi(locationService, accuWeatherService, openWeatherService),
		listCleanerTasks: initiateTaskList(&cfg.ExpirationDuration, logger, redisWeatherStorage, redisLocationStorage),
	}
}

func initiateTaskList(cfg *config.ExpirationDuration, logger log.Logger, weatherStorage storage.WeatherStorage, locationStorage storage.LocationStorage) []task.Task {
	return []task.Task{
		task.NewCoordsCleaner(logger, cfg, locationStorage),
		task.NewCurrentWeatherCleaner(logger, cfg, weatherStorage),
		task.NewHourlyWeatherCleaner(logger, cfg, weatherStorage),
		task.NewDailyWeatherCleaner(logger, cfg, weatherStorage),
	}
}

func (s *Server) Start(ctx context.Context) {
	s.wg.Add(2)

	go s.addHandleGraphqlToServer()
	go s.workflowForTask(ctx)

	s.wg.Wait()
}

func (s *Server) workflowForTask(ctx context.Context) {
	defer s.wg.Done()

	s.logger.Info("start run task cleaners")

	wgForTasks := &sync.WaitGroup{}

	for {
		time.Sleep(s.cfg.ExpirationDuration.TaskPeriod)

		for _, cleaner := range s.listCleanerTasks {
			s.runTask(ctx, cleaner, wgForTasks)
		}

		wgForTasks.Wait()
	}
}

func (s *Server) runTask(ctx context.Context, task task.Task, wgForTasks *sync.WaitGroup) {
	wgForTasks.Add(1)
	go func() {
		defer wgForTasks.Done()

		task.Run(ctx)
	}()
}

func (s *Server) addHandleGraphqlToServer() {
	defer s.wg.Done()
	endpoint := "/graphql/query"

	s.logger.Info("start add graphql handler to http server")

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(
		graph.Config{
			Resolvers: &graph.Resolver{
				GraphqlApi: s.graphqlApi,
			},
		},
	))

	http2.Handle(endpoint, srv)

	if s.cfg.Server.PlaygroundEnable {
		http2.Handle("/", playground.Handler("GraphQL playground", endpoint))
	}

	s.logger.Fatal(http2.ListenAndServe(fmt.Sprintf(":%d", s.cfg.Server.Port), nil))
}
