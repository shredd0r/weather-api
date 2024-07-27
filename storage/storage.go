package storage

import (
	"context"
	"errors"
	"weather-api/dto"
)

var (
	ErrNotFound = errors.New("volume not found")
)

type WeatherStorage interface {
	GetCurrentWeatherBy(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (*dto.CurrentWeather, error)
	GetHourlyWeatherBy(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (*[]*dto.HourlyWeather, error)
	GetDailyWeatherBy(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (*[]*dto.DailyWeather, error)
	SaveCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather *dto.CurrentWeather) error
	SaveHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather *[]*dto.HourlyWeather) error
	SaveDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather *[]*dto.DailyWeather) error
}

type LocationStorage interface {
	AddNewCoords(ctx context.Context, coords *dto.Coords, addressHash string) error
	LastTimeGetAddressHash(ctx context.Context, coords *dto.Coords) (uint64, error)
	UpdateLastTimeGetAddressHash(ctx context.Context, coords *dto.Coords) error
	GetAddressHashByCoords(ctx context.Context, coords *dto.Coords) (string, error)
	GetLocation(ctx context.Context, addressHash string) (*dto.Location, error)
	SaveLocation(ctx context.Context, location *dto.Location) error
}
