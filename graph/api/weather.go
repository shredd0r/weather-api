package api

import (
	"context"
	"github.com/shredd0r/weather-api/dto"
	"github.com/shredd0r/weather-api/graph/model"
	"github.com/shredd0r/weather-api/service"
)

type WeatherGraphqlApi struct {
	locationService    service.LocationService
	accuWeatherService service.WeatherService
	openWeatherService service.WeatherService
}

func NewWeatherGraphqlApi(locationService service.LocationService,
	accuWeatherService service.WeatherService,
	openWeatherService service.WeatherService) *WeatherGraphqlApi {
	return &WeatherGraphqlApi{
		accuWeatherService: accuWeatherService,
		openWeatherService: openWeatherService,
	}
}

func (p *WeatherGraphqlApi) FindGeocoding(ctx context.Context, input *dto.GeocodingRequest) (*[]*dto.Geocoding, error) {
	return p.locationService.FindGeocoding(ctx, input)
}

func (p *WeatherGraphqlApi) CurrentWeather(ctx context.Context, input *model.WeatherRequest) (*dto.CurrentWeather, error) {
	request := p.mapToWeatherRequest(input)
	switch input.Forecaster {
	case dto.WeatherForecasterAccuWeather:
		return p.accuWeatherService.CurrentWeather(ctx, request)
	case dto.WeatherForecasterOpenWeather:
		return p.openWeatherService.CurrentWeather(ctx, request)
	}
	return nil, nil
}

func (p *WeatherGraphqlApi) HourlyWeather(ctx context.Context, input *model.WeatherRequest) (*[]*dto.HourlyWeather, error) {
	request := p.mapToWeatherRequest(input)
	switch input.Forecaster {
	case dto.WeatherForecasterAccuWeather:
		return p.accuWeatherService.HourlyWeather(ctx, request)
	case dto.WeatherForecasterOpenWeather:
		return p.openWeatherService.HourlyWeather(ctx, request)
	}
	return nil, nil
}

func (p *WeatherGraphqlApi) DailyWeather(ctx context.Context, input *model.WeatherRequest) (*[]*dto.DailyWeather, error) {
	request := p.mapToWeatherRequest(input)
	switch input.Forecaster {
	case dto.WeatherForecasterAccuWeather:
		return p.accuWeatherService.DailyWeather(ctx, request)
	case dto.WeatherForecasterOpenWeather:
		return p.openWeatherService.DailyWeather(ctx, request)
	}
	return nil, nil
}

func (p *WeatherGraphqlApi) mapToWeatherRequest(input *model.WeatherRequest) dto.WeatherRequest {
	return dto.WeatherRequest{
		Coords: &dto.Coords{
			Latitude:  input.Coords.Latitude,
			Longitude: input.Coords.Longitude,
		},
		Locale: input.Locale,
		Unit:   input.Unit,
	}
}
