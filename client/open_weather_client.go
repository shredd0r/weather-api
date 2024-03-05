package client

import (
	"log"
	"net/http"
	"weather-api/dto"
)

type OpenWeatherClient struct {
	client *Client
}

func NewOpenWeatherClient(logger *log.Logger, httpClient *http.Client) *AccuWeatherClient {
	return &AccuWeatherClient{
		client: NewClient("https://api.openweathermap.org/", logger, httpClient),
	}
}

func (client *Client) GetCurrentWeatherInfo(request dto.OpenWeatherRequestDto) (*dto.OpenWeatherCurrentWeatherResponseDto, error) {
	return nil, nil
}

func (client *Client) GetForecastWeatherInfo(request dto.OpenWeatherRequestDto) (*dto.OpenWeatherHourlyWeatherResponseDto, error) {
	return nil, nil
}
