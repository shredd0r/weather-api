package task

import (
	"context"
	"github.com/shredd0r/weather-api/config"
	"github.com/shredd0r/weather-api/log"
	"github.com/shredd0r/weather-api/storage"
)

type DailyWeatherCleaner struct {
	*baseWeatherCleaner
}

func NewDailyWeatherCleaner(logger log.Logger, cfg *config.ExpirationDuration, storage storage.WeatherStorage) Task {
	return &DailyWeatherCleaner{
		baseWeatherCleaner: newBaseWeatherCleaner(logger, cfg, storage),
	}
}

func (c *DailyWeatherCleaner) Run(ctx context.Context) {
	c.workFlowForRemoveWeatherForAllForecasterFromCache(ctx, c.storage.GetAllLastTimeUpdatedDailyWeather, c.storage.RemoveDailyWeather, c.storage.RemoveLastTimeUpdatedDailyWeather)
}
