package task

import (
	"context"
	"github.com/shredd0r/weather-api/config"
	"github.com/shredd0r/weather-api/log"
	"github.com/shredd0r/weather-api/storage"
)

type HourlyWeatherCleaner struct {
	*baseWeatherCleaner
}

func NewHourlyWeatherCleaner(logger log.Logger, cfg *config.ExpirationDuration, storage storage.WeatherStorage) Task {
	return &HourlyWeatherCleaner{
		baseWeatherCleaner: newBaseWeatherCleaner(logger, cfg, storage),
	}
}

func (c *HourlyWeatherCleaner) Run(ctx context.Context) {
	c.workFlowForRemoveWeatherForAllForecasterFromCache(ctx, c.storage.GetAllLastTimeUpdatedHourlyWeather, c.storage.RemoveHourlyWeather, c.storage.RemoveLastTimeUpdatedHourlyWeather)
}
