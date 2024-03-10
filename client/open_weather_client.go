package client

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"weather-api/dto"
)

type OpenWeatherClient struct {
	*Client
}

func NewOpenWeatherClient(logger *logrus.Logger, httpClient *http.Client) *AccuWeatherClient {
	return &AccuWeatherClient{
		NewClient("https://api.openweathermap.org/", logger, httpClient),
	}
}

func (c *Client) GetCurrentWeatherInfo(request dto.OpenWeatherRequestDto) (*dto.OpenWeatherCurrentWeatherResponseDto, error) {
	return nil, nil
}

func (c *Client) GetForecastWeatherInfo(request dto.OpenWeatherRequestDto) (*dto.OpenWeatherHourlyWeatherResponseDto, error) {
	return nil, nil
}
