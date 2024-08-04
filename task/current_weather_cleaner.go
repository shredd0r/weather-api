package task

import (
	"context"
	"github.com/shredd0r/weather-api/config"
	"github.com/shredd0r/weather-api/log"
	"github.com/shredd0r/weather-api/storage"
)

type CurrentWeatherCleaner struct {
	*baseWeatherCleaner
}

func NewCurrentWeatherCleaner(logger log.Logger, cfg *config.ExpirationDuration, storage storage.WeatherStorage) Task {
	return &CurrentWeatherCleaner{
		baseWeatherCleaner: newBaseWeatherCleaner(logger, cfg, storage),
	}
}

func (c *CurrentWeatherCleaner) Run(ctx context.Context) {
	c.workFlowForRemoveWeatherForAllForecasterFromCache(ctx, c.storage.GetAllLastTimeUpdatedCurrentWeather, c.storage.RemoveCurrentWeather, c.storage.RemoveLastTimeUpdatedCurrentWeather)
}
