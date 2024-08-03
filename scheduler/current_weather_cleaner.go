package scheduler

import (
	"context"
	"sync"
	"weather-api/config"
	"weather-api/log"
	"weather-api/storage"
)

type CurrentWeatherCleaner struct {
	*baseWeatherCleaner
}

func NewCurrentWeatherCleaner(logger log.Logger, cfg *config.ExpirationDuration, wg *sync.WaitGroup, storage storage.WeatherStorage) Scheduler {
	return &CurrentWeatherCleaner{
		baseWeatherCleaner: newBaseWeatherCleaner(logger, cfg, wg, storage),
	}
}

func (c *CurrentWeatherCleaner) Run(ctx context.Context) {
	c.workFlowForRemoveWeatherForAllForecasterFromCache(ctx, c.storage.GetAllLastTimeUpdatedCurrentWeather, c.storage.RemoveCurrentWeather, c.storage.RemoveLastTimeUpdatedCurrentWeather)
}
