package scheduler

import (
	"context"
	"sync"
	"time"
	"weather-api/config"
	"weather-api/dto"
	"weather-api/log"
	"weather-api/storage"
)

type baseWeatherCleaner struct {
	logger  log.Logger
	cfg     *config.ExpirationDuration
	wg      *sync.WaitGroup
	storage storage.WeatherStorage
}

func newBaseWeatherCleaner(logger log.Logger, cfg *config.ExpirationDuration, wg *sync.WaitGroup, storage storage.WeatherStorage) *baseWeatherCleaner {
	return &baseWeatherCleaner{
		logger:  logger,
		cfg:     cfg,
		wg:      wg,
		storage: storage,
	}
}

func (c *baseWeatherCleaner) workFlowForRemoveWeatherFromCache(
	ctx context.Context,
	forecaster dto.WeatherForecaster,
	methodForGetAllLastTimeUpdated func(context.Context, dto.WeatherForecaster) (map[string]int64, error),
	methodForRemoveWeather func(context.Context, string, dto.WeatherForecaster) error,
	methodForRemoveLastTimeUpdated func(context.Context, string, dto.WeatherForecaster) error) {
	go func() {
		defer c.wg.Done()

		mapWithAddressHashAndLastTimeUpdated, err := methodForGetAllLastTimeUpdated(ctx, forecaster)
		if err != nil {
			c.logger.Error("error when try get all last time updated from cache")
			return
		}

		now := time.Now().UnixMilli()
		for addressHash, lastTimeUpdated := range mapWithAddressHashAndLastTimeUpdated {
			if isRowNeedRemove(now, lastTimeUpdated, c.cfg.WeatherInfo) {
				c.removeDataAboutWeatherFromCache(ctx, addressHash, forecaster, methodForRemoveWeather, methodForRemoveLastTimeUpdated)
			}
		}

		c.wg.Wait()
	}()
}

func (c *baseWeatherCleaner) workFlowForRemoveWeatherForAllForecasterFromCache(
	ctx context.Context,
	methodForGetAllLastTimeUpdated func(context.Context, dto.WeatherForecaster) (map[string]int64, error),
	methodForRemoveWeather func(context.Context, string, dto.WeatherForecaster) error,
	methodForRemoveLastTimeUpdated func(context.Context, string, dto.WeatherForecaster) error) {
	c.wg.Add(2)

	c.workFlowForRemoveWeatherFromCache(ctx, dto.WeatherForecasterAccuWeather, methodForGetAllLastTimeUpdated, methodForRemoveWeather, methodForRemoveLastTimeUpdated)
	c.workFlowForRemoveWeatherFromCache(ctx, dto.WeatherForecasterOpenWeather, methodForGetAllLastTimeUpdated, methodForRemoveWeather, methodForRemoveLastTimeUpdated)

	c.wg.Wait()
}

func (c *baseWeatherCleaner) removeDataAboutWeatherFromCache(
	ctx context.Context,
	addressHash string,
	forecaster dto.WeatherForecaster,
	methodForRemoveWeather func(context.Context, string, dto.WeatherForecaster) error,
	methodForRemoveLastTimeUpdated func(context.Context, string, dto.WeatherForecaster) error,
) {
	c.wg.Add(2)

	go func() {
		defer c.wg.Done()

		err := methodForRemoveWeather(ctx, addressHash, forecaster)
		if err != nil {
			c.logger.Error("error when try remove weather from cache")
		}
	}()

	go func() {
		defer c.wg.Done()

		err := methodForRemoveLastTimeUpdated(ctx, addressHash, forecaster)
		if err != nil {
			c.logger.Error("error when try remove last time updated weather from cache")
		}
	}()
}
