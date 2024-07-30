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
	SaveCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather dto.CurrentWeather) error
	SaveHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather []*dto.HourlyWeather) error
	SaveDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather []*dto.DailyWeather) error
	SaveUpdatedTimeCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, lastTime int64) error
	SaveUpdatedTimeHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, lastTime int64) error
	SaveUpdatedTimeDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, lastTime int64) error
	GetLastTimeUpdatedCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (int64, error)
	GetLastTimeUpdatedHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (int64, error)
	GetLastTimeUpdatedDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (int64, error)
}

type LocationStorage interface {
	AddNewCoords(ctx context.Context, coords *dto.Coords, addressHash string) error
	GetLastTimeGetAddressHash(ctx context.Context, coords *dto.Coords) (int64, error)
	UpdateLastTimeGetAddressHash(ctx context.Context, coords *dto.Coords, time int64) error
	GetAddressHashByCoords(ctx context.Context, coords *dto.Coords) (string, error)
	GetLocation(ctx context.Context, addressHash string) (*dto.Location, error)
	SaveLocation(ctx context.Context, location dto.Location, addressHash string) error
}
