package provider

import (
	"context"
	"github.com/sirupsen/logrus"
	"weather-api/client/http"
	"weather-api/dto"
)

type OpenWeatherProvider struct {
	apiKey string
	logger *logrus.Logger
	client http.OpenWeatherInterface
}

func (p OpenWeatherProvider) CurrentWeather(ctx context.Context, request dto.WeatherRequestProviderDto) (*dto.CurrentWeather, error) {
	panic("implement me")

}

func (p OpenWeatherProvider) HourlyWeather(ctx context.Context, request dto.WeatherRequestProviderDto) (*[]*dto.HourlyWeather, error) {
	//TODO implement me
	panic("implement me")
}

func (p OpenWeatherProvider) DailyWeather(ctx context.Context, request dto.WeatherRequestProviderDto) (*[]*dto.DailyWeather, error) {
	//TODO implement me
	panic("implement me")
}
