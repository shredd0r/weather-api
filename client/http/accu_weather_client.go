package http

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
	"weather-api/dto"
)

type AccuWeatherInterface interface {
	GetCurrentWeatherInfo(request dto.AccuWeatherRequestDto) (*dto.AccuWeatherCurrentResponseDto, error)
	GetHourlyWeatherInfo(request dto.AccuWeatherRequestDto) (*[]dto.AccuWeatherHourlyResponseDto, error)
	GetDailyWeatherInfo(request dto.AccuWeatherRequestDto) (*dto.AccuWeatherDailyResponseDto, error)
	GetGeoPositionSearch(request dto.AccuWeatherGeoPositionRequestDto) (*dto.AccuWeatherGeoPositionResponseDto, error)
}

type AccuWeatherClient struct {
	client *Client
}

func NewAccuWeatherClient(log *logrus.Logger, httpClient *http.Client, apiKey string) *AccuWeatherClient {
	return &AccuWeatherClient{
		NewClient("http://dataservice.accuweather.com/", log, httpClient, apiKey),
	}
}

func (c *AccuWeatherClient) GetCurrentWeatherInfo(request dto.AccuWeatherRequestDto) (*dto.AccuWeatherCurrentResponseDto, error) {
	var urlForRequest = c.client.BaseURL + fmt.Sprintf("currentconditions/v1/%s?", request.LocationKey) + c.getQueryParamsForRequest(request)
	response, err := HttpGetAndGetResponse[[]*dto.AccuWeatherCurrentResponseDto](c.client.httpClient, c.client.log, GetHttpRequestBy(urlForRequest))

	if err == nil {
		return (*response)[0], err
	}

	return nil, err
}

func (c *AccuWeatherClient) GetHourlyWeatherInfo(request dto.AccuWeatherRequestDto) (*[]dto.AccuWeatherHourlyResponseDto, error) {
	var urlForRequest = c.client.BaseURL + fmt.Sprintf("forecasts/v1/hourly/12hour/%s?", request.LocationKey) + c.getQueryParamsForRequest(request)
	return HttpGetAndGetResponse[[]dto.AccuWeatherHourlyResponseDto](c.client.httpClient, c.client.log, GetHttpRequestBy(urlForRequest))
}

func (c *AccuWeatherClient) GetDailyWeatherInfo(request dto.AccuWeatherRequestDto) (*dto.AccuWeatherDailyResponseDto, error) {
	var urlForRequest = c.client.BaseURL + fmt.Sprintf("forecasts/v1/daily/5day/%s?", request.LocationKey) + c.getQueryParamsForRequest(request)
	return HttpGetAndGetResponse[dto.AccuWeatherDailyResponseDto](c.client.httpClient, c.client.log, GetHttpRequestBy(urlForRequest))

}

func (c *AccuWeatherClient) GetGeoPositionSearch(request dto.AccuWeatherGeoPositionRequestDto) (*dto.AccuWeatherGeoPositionResponseDto, error) {
	var urlForRequest = c.client.BaseURL + "locations/v1/cities/geoposition/search?" + c.getQueryParamsForGeoPosition(request)
	return HttpGetAndGetResponse[dto.AccuWeatherGeoPositionResponseDto](c.client.httpClient, c.client.log, GetHttpRequestBy(urlForRequest))
}

func (c *AccuWeatherClient) getQueryParamsForBase(request dto.AccuWeatherBaseRequestDto) url.Values {
	queryParams := url.Values{
		"apikey":   {c.client.apiKey},
		"language": {request.Language},
		"details":  {strconv.FormatBool(request.Details)},
		"metric":   {strconv.FormatBool(request.Metric)},
	}
	return queryParams
}

func (c *AccuWeatherClient) getQueryParamsForRequest(request dto.AccuWeatherRequestDto) string {
	queryParams := c.getQueryParamsForBase(request.AccuWeatherBaseRequestDto)
	queryParams.Add("locationKey", request.LocationKey)

	return queryParams.Encode()
}

func (c *AccuWeatherClient) getQueryParamsForGeoPosition(request dto.AccuWeatherGeoPositionRequestDto) string {
	q := fmt.Sprintf("%f,%f", request.Latitude, request.Longitude)
	queryParams := url.Values{
		"apikey": {c.client.apiKey},
		"q":      {q},
	}

	return queryParams.Encode()
}
