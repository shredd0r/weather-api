package task

import (
	"context"
	"github.com/shredd0r/weather-api/config"
	"github.com/shredd0r/weather-api/dto"
	"github.com/shredd0r/weather-api/log"
	"github.com/shredd0r/weather-api/storage"
	"sync"
	"time"
)

type baseWeatherCleaner struct {
	logger  log.Logger
	cfg     *config.ExpirationDuration
	wg      *sync.WaitGroup
	storage storage.WeatherStorage
}

func newBaseWeatherCleaner(logger log.Logger, cfg *config.ExpirationDuration, storage storage.WeatherStorage) *baseWeatherCleaner {
	return &baseWeatherCleaner{
		logger:  logger,
		cfg:     cfg,
		wg:      &sync.WaitGroup{},
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
		c.logger.Debugf("start remove weather from cache for forecaster: %s", forecaster)
		defer c.wg.Done()

		mapWithAddressHashAndLastTimeUpdated, err := methodForGetAllLastTimeUpdated(ctx, forecaster)
		if err != nil {
			c.logger.Error("error when try get all last time updated from cache")
			return
		}

		c.logger.Debugf("fined %d weather from cache", len(mapWithAddressHashAndLastTimeUpdated))

		now := time.Now().UnixMilli()
		for addressHash, lastTimeUpdated := range mapWithAddressHashAndLastTimeUpdated {
			if isRowNeedRemove(now, lastTimeUpdated, c.cfg.WeatherInfo) {
				c.logger.Debugf("weather for %s with addressHash: %s need remove", forecaster, addressHash)
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

		c.logger.Debug("weather removed from cache")
	}()

	go func() {
		defer c.wg.Done()

		err := methodForRemoveLastTimeUpdated(ctx, addressHash, forecaster)
		if err != nil {
			c.logger.Error("error when try remove last time updated weather from cache")
		}
		c.logger.Debug("last time updated weather removed from cache")
	}()
}
