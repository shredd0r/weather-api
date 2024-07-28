package storage

import (
	"context"
	"weather-api/dto"
)

type RedisWeatherStorage struct {
}

func (r RedisWeatherStorage) GetCurrentWeatherBy(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (*dto.CurrentWeather, error) {
	//TODO implement me
	return nil, ErrNotFound
}

func (r RedisWeatherStorage) GetHourlyWeatherBy(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (*[]*dto.HourlyWeather, error) {
	//TODO implement me
	return nil, ErrNotFound

}

func (r RedisWeatherStorage) GetDailyWeatherBy(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (*[]*dto.DailyWeather, error) {
	//TODO implement me
	return nil, ErrNotFound
}

func (r RedisWeatherStorage) SaveCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather dto.CurrentWeather) error {
	//TODO implement me
	return ErrNotFound
}

func (r RedisWeatherStorage) SaveHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather []*dto.HourlyWeather) error {
	//TODO implement me
	return ErrNotFound
}

func (r RedisWeatherStorage) SaveDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather []*dto.DailyWeather) error {
	//TODO implement me
	return ErrNotFound
}
