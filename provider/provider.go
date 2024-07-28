package provider

import (
	"context"
	"weather-api/dto"
)

type WeatherProvider interface {
	CurrentWeather(ctx context.Context, request *dto.WeatherRequestProviderDto) (*dto.CurrentWeather, error)
	HourlyWeather(ctx context.Context, request *dto.WeatherRequestProviderDto) (*[]*dto.HourlyWeather, error)
	DailyWeather(ctx context.Context, request *dto.WeatherRequestProviderDto) (*[]*dto.DailyWeather, error)
}

type CacheWeatherProvider interface {
	CurrentWeather(ctx context.Context, request *dto.WeatherRequestProviderDto, forecaster dto.WeatherForecaster) (*dto.CurrentWeather, error)
	HourlyWeather(ctx context.Context, request *dto.WeatherRequestProviderDto, forecaster dto.WeatherForecaster) (*[]*dto.HourlyWeather, error)
	DailyWeather(ctx context.Context, request *dto.WeatherRequestProviderDto, forecaster dto.WeatherForecaster) (*[]*dto.DailyWeather, error)
}

type LocationProvider interface {
	GetAddressHashByCoords(ctx context.Context, coords *dto.Coords) (string, error)
	LocationByAddressHash(ctx context.Context, coords *dto.Coords, addressHash string) (*dto.Location, error)
}
