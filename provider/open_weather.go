package provider

import (
	"context"

	"github.com/shredd0r/weather-api/client/http"
	"github.com/shredd0r/weather-api/dto"
	"github.com/shredd0r/weather-api/log"
	"github.com/shredd0r/weather-api/mapper"
)

type OpenWeatherProvider struct {
	logger log.Logger
	client http.OpenWeatherInterface
}

func NewOpenWeatherProvider(logger log.Logger, client http.OpenWeatherInterface) WeatherProvider {
	return &OpenWeatherProvider{
		logger: logger,
		client: client,
	}
}

func (p OpenWeatherProvider) CurrentWeather(ctx context.Context, request *dto.WeatherRequestProvider) (*dto.CurrentWeather, error) {
	resp, err := p.client.GetCurrentWeatherInfo(mapper.OpenWeatherGetRequest(request))

	if err != nil {
		return nil, err
	}

	return mapper.OpenWeatherMapCurrentWeather(resp), nil
}

func (p OpenWeatherProvider) HourlyWeather(ctx context.Context, request *dto.WeatherRequestProvider) (*[]*dto.HourlyWeather, error) {
	resp, err := p.client.GetForecastWeatherInfo(mapper.OpenWeatherGetForecastRequest(request))

	if err != nil {
		return nil, err
	}

	return mapper.OpenWeatherMapHourlyWeahters(resp), nil
}

func (p OpenWeatherProvider) DailyWeather(ctx context.Context, request *dto.WeatherRequestProvider) (*[]*dto.DailyWeather, error) {
	resp, err := p.client.GetForecastWeatherInfo(mapper.OpenWeatherGetForecastRequest(request))

	if err != nil {
		return nil, err
	}

	return mapper.OpenWeatherMapDailyWeahters(resp), nil
}

func (p OpenWeatherProvider) getUnits(unit dto.Unit) dto.OpenWeatherUnits {
	switch unit {
	case dto.UnitMetric:
		return dto.OpenWeatherUnitsMetric
	case dto.UnitImperial:
		return dto.OpenWeatherUnitsImperial
	default:
		return dto.OpenWeatherUnitsStandard
	}
}
func (p OpenWeatherProvider) getForecastRequest(request *dto.WeatherRequestProvider) dto.OpenWeatherForecastRequestDto {
	return dto.OpenWeatherForecastRequestDto{
		OpenWeatherRequestDto: p.getRequest(request),
	}
}

func (p OpenWeatherProvider) getRequest(request *dto.WeatherRequestProvider) dto.OpenWeatherRequestDto {
	return dto.OpenWeatherRequestDto{
		Latitude:  request.Location.Coords.Latitude,
		Longitude: request.Location.Coords.Longitude,
		Language:  request.Locale,
		Units:     p.getUnits(request.Unit),
	}
}
