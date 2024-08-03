package scheduler

import (
	"context"
	"sync"
	"weather-api/config"
	"weather-api/log"
	"weather-api/storage"
)

type HourlyWeatherCleaner struct {
	*baseWeatherCleaner
}

func NewHourlyWeatherCleaner(logger log.Logger, cfg *config.ExpirationDuration, wg *sync.WaitGroup, storage storage.WeatherStorage) Scheduler {
	return &HourlyWeatherCleaner{
		baseWeatherCleaner: newBaseWeatherCleaner(logger, cfg, wg, storage),
	}
}

func (c *HourlyWeatherCleaner) Run(ctx context.Context) {
	c.workFlowForRemoveWeatherForAllForecasterFromCache(ctx, c.storage.GetAllLastTimeUpdatedHourlyWeather, c.storage.RemoveHourlyWeather, c.storage.RemoveLastTimeUpdatedHourlyWeather)
}
