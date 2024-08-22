//go:generate mockgen -source=storage.go -destination=mock/storage.go
package storage

import (
	"context"
	"errors"

	"github.com/shredd0r/weather-api/dto"
)

var (
	ErrNotFound = errors.New("volume not found")
)

type WeatherStorage interface {
	SaveCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather *dto.CurrentWeather) error
	SaveHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather *[]*dto.HourlyWeather) error
	SaveDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather *[]*dto.DailyWeather) error

	GetCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (*dto.CurrentWeather, error)
	GetHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (*[]*dto.HourlyWeather, error)
	GetDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (*[]*dto.DailyWeather, error)

	RemoveCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) error
	RemoveHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) error
	RemoveDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) error

	SaveUpdatedTimeCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, lastTime int64) error
	SaveUpdatedTimeHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, lastTime int64) error
	SaveUpdatedTimeDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, lastTime int64) error

	GetAllLastTimeUpdatedCurrentWeather(ctx context.Context, forecaster dto.WeatherForecaster) (map[string]int64, error)
	GetAllLastTimeUpdatedHourlyWeather(ctx context.Context, forecaster dto.WeatherForecaster) (map[string]int64, error)
	GetAllLastTimeUpdatedDailyWeather(ctx context.Context, forecaster dto.WeatherForecaster) (map[string]int64, error)

	RemoveLastTimeUpdatedCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) error
	RemoveLastTimeUpdatedHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) error
	RemoveLastTimeUpdatedDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) error
}

type LocationStorage interface {
	GetLocation(ctx context.Context, addressHash string) (*dto.Location, error)
	SaveLocation(ctx context.Context, location *dto.Location, addressHash string) error

	AddCoords(ctx context.Context, coords *dto.Coords, addressHash string) error
	RemoveCoords(ctx context.Context, coords *dto.Coords) error

	UpdateLastTimeUseCoords(ctx context.Context, coords *dto.Coords, time int64) error
	GetAllLastTimeUseCoords(ctx context.Context) (map[dto.Coords]int64, error)

	GetAddressHashByCoords(ctx context.Context, coords *dto.Coords) (string, error)
}
