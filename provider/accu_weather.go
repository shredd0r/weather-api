package provider

import (
	"context"
	"github.com/sirupsen/logrus"
	"weather-api/client/http"
	"weather-api/dto"
)

type AccuWeatherProvider struct {
	logger *logrus.Logger
	client http.AccuWeatherInterface
}

func NewAccuWeatherProvider(logger *logrus.Logger, client http.AccuWeatherInterface) WeatherProvider {
	return &AccuWeatherProvider{
		logger: logger,
		client: client,
	}
}

func (p AccuWeatherProvider) CurrentWeather(ctx context.Context, request dto.WeatherRequestProviderDto) (*dto.CurrentWeather, error) {
	resp, err := p.client.GetCurrentWeatherInfo(dto.AccuWeatherRequestDto{
		AccuWeatherBaseRequestDto: dto.AccuWeatherBaseRequestDto{
			Language: "uk-UA",
			Metric:   request.Unit == dto.UnitMetric,
			Details:  true,
		},
		LocationKey: request.Location.AccuWeatherLocationKey,
	})

	if err != nil {
		return nil, err
	}

	return &dto.CurrentWeather{
		EpochTime:            resp.EpochTime,
		Visibility:           resp.Visibility.Metric.Value,
		CurrentTemperature:   resp.Temperature.Metric.Value,
		FeelsLikeTemperature: resp.RealFeelTemperature.Metric.Value,
		MobileLink:           resp.MobileLink,
		Link:                 resp.Link,
	}, nil
}

func (p AccuWeatherProvider) HourlyWeather(ctx context.Context, request dto.WeatherRequestProviderDto) (*[]*dto.HourlyWeather, error) {
	//TODO implement me
	panic("implement me")
}

func (p AccuWeatherProvider) DailyWeather(ctx context.Context, request dto.WeatherRequestProviderDto) (*[]*dto.DailyWeather, error) {
	//TODO implement me
	panic("implement me")
}
