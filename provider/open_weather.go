package provider

import (
	"context"
	"github.com/sirupsen/logrus"
	"weather-api/client/http"
	"weather-api/dto"
	"weather-api/util"
)

type OpenWeatherProvider struct {
	logger *logrus.Logger
	client http.OpenWeatherInterface
}

func (p OpenWeatherProvider) CurrentWeather(ctx context.Context, request dto.WeatherRequestProviderDto) (*dto.CurrentWeather, error) {
	resp, err := p.client.GetCurrentWeatherInfo(dto.OpenWeatherRequestDto{
		Latitude:  request.Location.Coords.Latitude,
		Longitude: request.Location.Coords.Longitude,
		Language:  request.Locale,
		Units:     p.getUnits(request.Unit),
	})

	if err != nil {
		return nil, err
	}

	return &dto.CurrentWeather{
		EpochTime:            resp.EpochTime,
		Visibility:           util.PercentToFloat64Pointer(&resp.Visibility),
		CurrentTemperature:   &resp.Main.Temperature,
		MinTemperature:       &resp.Main.MinTemperature,
		MaxTemperature:       &resp.Main.MaxTemperature,
		FeelsLikeTemperature: &resp.Main.FeelsLike,
		IconResource:         nil,
	}, nil
}

func (p OpenWeatherProvider) HourlyWeather(ctx context.Context, request dto.WeatherRequestProviderDto) (*[]*dto.HourlyWeather, error) {
	resp, err := p.client.GetForecastWeatherInfo(dto.OpenWeatherForecastRequestDto{
		OpenWeatherRequestDto: dto.OpenWeatherRequestDto{
			Latitude:  request.Location.Coords.Latitude,
			Longitude: request.Location.Coords.Longitude,
			Language:  request.Locale,
			Units:     p.getUnits(request.Unit),
		},
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (p OpenWeatherProvider) DailyWeather(ctx context.Context, request dto.WeatherRequestProviderDto) (*[]*dto.DailyWeather, error) {
	//TODO implement me
	panic("implement me")
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
