package scheduler

import (
	"context"
	"sync"
	"weather-api/config"
	"weather-api/log"
	"weather-api/storage"
)

type DailyWeatherCleaner struct {
	*baseWeatherCleaner
}

func NewDailyWeatherCleaner(logger log.Logger, cfg *config.ExpirationDuration, wg *sync.WaitGroup, storage storage.WeatherStorage) Scheduler {
	return &DailyWeatherCleaner{
		baseWeatherCleaner: newBaseWeatherCleaner(logger, cfg, wg, storage),
	}
}

func (c *DailyWeatherCleaner) Run(ctx context.Context) {
	c.workFlowForRemoveWeatherForAllForecasterFromCache(ctx, c.storage.GetAllLastTimeUpdatedDailyWeather, c.storage.RemoveDailyWeather, c.storage.RemoveLastTimeUpdatedDailyWeather)
}
