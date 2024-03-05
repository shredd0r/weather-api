package client

import (
	"log"
	"net/http"
	"weather-api/dto"
)

type AccuWeatherClient struct {
	client *Client
}

func NewAccuWeatherClient(logger *log.Logger, httpClient *http.Client) *AccuWeatherClient {
	return &AccuWeatherClient{
		client: NewClient("http://dataservice.accuweather.com/", logger, httpClient),
	}
}

func (client *AccuWeatherClient) GetCurrentWeatherInfo(request dto.AccuWeatherRequestDto) (*dto.AccuWeatherDailyResponseDto, error) {
	return nil, nil
}

func (client *AccuWeatherClient) GetHourlyWeatherInfo(request dto.AccuWeatherRequestDto) (*dto.AccuWeatherHourlyResponseDto, error) {
	return nil, nil
}

func (client *AccuWeatherClient) GetDailyWeatherInfo(request dto.AccuWeatherRequestDto) (*dto.AccuWeatherDailyResponseDto, error) {
	return nil, nil
}

func (client *AccuWeatherClient) GetGeoPositionSearch(request dto.AccuWeatherGeoPositionRequestDto) (*dto.AccuWeatherGeoPositionResponseDto, error) {
	return nil, nil
}
