package provider

import (
	"context"

	"github.com/shredd0r/weather-api/client/http"
	"github.com/shredd0r/weather-api/dto"
	"github.com/shredd0r/weather-api/log"
	"github.com/shredd0r/weather-api/mapper"
)

type AccuWeatherProvider struct {
	logger log.Logger
	client http.AccuWeatherInterface
}

func NewAccuWeatherProvider(logger log.Logger, client http.AccuWeatherInterface) WeatherProvider {
	return &AccuWeatherProvider{
		logger: logger,
		client: client,
	}
}

func (p AccuWeatherProvider) CurrentWeather(ctx context.Context, request *dto.WeatherRequestProvider) (*dto.CurrentWeather, error) {
	resp, err := p.client.GetCurrentWeatherInfo(mapper.AccuWeatherMapRequestForClient(request))

	if err != nil {
		return nil, err
	}

	return mapper.AccuWeatherMapToCurrentWeatherDto(request.Unit, resp), nil
}

func (p AccuWeatherProvider) HourlyWeather(ctx context.Context, request *dto.WeatherRequestProvider) (*[]*dto.HourlyWeather, error) {
	resp, err := p.client.GetHourlyWeatherInfo(mapper.AccuWeatherMapRequestForClient(request))

	if err != nil {
		return nil, err
	}

	return mapper.AccuWeatherMapToHourlyWeatherDtos(*resp), nil
}

func (p AccuWeatherProvider) DailyWeather(ctx context.Context, request *dto.WeatherRequestProvider) (*[]*dto.DailyWeather, error) {
	resp, err := p.client.GetDailyWeatherInfo(mapper.AccuWeatherMapRequestForClient(request))

	if err != nil {
		return nil, err
	}

	return mapper.AccuWeatherMapToDailyWeatherDtos(resp), err
}
