package provider

import (
	"context"
	"github.com/shredd0r/weather-api/dto"
)

type WeatherProvider interface {
	CurrentWeather(ctx context.Context, request *dto.WeatherRequestProvider) (*dto.CurrentWeather, error)
	HourlyWeather(ctx context.Context, request *dto.WeatherRequestProvider) (*[]*dto.HourlyWeather, error)
	DailyWeather(ctx context.Context, request *dto.WeatherRequestProvider) (*[]*dto.DailyWeather, error)
}

type CacheWeatherProvider interface {
	CurrentWeather(ctx context.Context, request *dto.WeatherRequestProvider, forecaster dto.WeatherForecaster) (*dto.CurrentWeather, error)
	HourlyWeather(ctx context.Context, request *dto.WeatherRequestProvider, forecaster dto.WeatherForecaster) (*[]*dto.HourlyWeather, error)
	DailyWeather(ctx context.Context, request *dto.WeatherRequestProvider, forecaster dto.WeatherForecaster) (*[]*dto.DailyWeather, error)
}

type LocationProvider interface {
	FindGeocoding(ctx context.Context, request *dto.GeocodingRequest) (*[]*dto.Geocoding, error)
	GetAddressHashByCoords(ctx context.Context, coords *dto.Coords) (string, error)
	LocationByAddressHash(ctx context.Context, coords *dto.Coords, addressHash string) (*dto.Location, error)
}
