package client

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
	"weather-api/config"
	"weather-api/dto"
)

var (
	StatusCodeNot200 = errors.New("api returned code != 200 OK")
)

type AccuWeatherClient struct {
	*Client
}

func NewAccuWeatherClient(log *logrus.Logger, httpClient *http.Client, cfg *config.WeatherApiKey) *AccuWeatherClient {
	return &AccuWeatherClient{
		NewClient("http://dataservice.accuweather.com/", log, httpClient, cfg),
	}
}

func (c *AccuWeatherClient) GetCurrentWeatherInfo(request dto.AccuWeatherRequestDto) (*dto.AccuWeatherCurrentResponseDto, error) {
	var urlForRequest = c.BaseURL + fmt.Sprintf("currentconditions/v1/%s?", request.LocationKey) + c.getQueryParamsForRequest(request)
	response, err := HttpGetAndGetResponse[[]*dto.AccuWeatherCurrentResponseDto](c.httpClient, c.log, urlForRequest)

	if err == nil {
		return (*response)[0], err
	}

	return nil, err
}

func (c *AccuWeatherClient) GetHourlyWeatherInfo(request dto.AccuWeatherRequestDto) (*[]*dto.AccuWeatherHourlyResponseDto, error) {
	var urlForRequest = c.BaseURL + fmt.Sprintf("forecasts/v1/hourly/12hour/%s?", request.LocationKey) + c.getQueryParamsForRequest(request)
	return HttpGetAndGetResponse[[]*dto.AccuWeatherHourlyResponseDto](c.httpClient, c.log, urlForRequest)
}

func (c *AccuWeatherClient) GetDailyWeatherInfo(request dto.AccuWeatherRequestDto) (*dto.AccuWeatherDailyResponseDto, error) {
	var urlForRequest = c.BaseURL + fmt.Sprintf("forecasts/v1/daily/5day/%s?", request.LocationKey) + c.getQueryParamsForRequest(request)
	return HttpGetAndGetResponse[dto.AccuWeatherDailyResponseDto](c.httpClient, c.log, urlForRequest)

}

func (c *AccuWeatherClient) GetGeoPositionSearch(request dto.AccuWeatherGeoPositionRequestDto) (*dto.AccuWeatherGeoPositionResponseDto, error) {
	var urlForRequest = c.BaseURL + "locations/v1/cities/geoposition/search?" + c.getQueryParamsForGeoPosition(request)
	return HttpGetAndGetResponse[dto.AccuWeatherGeoPositionResponseDto](c.httpClient, c.log, urlForRequest)
}

func (c *AccuWeatherClient) getQueryParamsForBase(request dto.AccuWeatherBaseRequestDto) url.Values {
	queryParams := url.Values{
		"apikey":   {c.apiKey},
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
	q := fmt.Sprintf("%.f,%.f", request.Latitude, request.Longitude)
	queryParams := c.getQueryParamsForBase(request.AccuWeatherBaseRequestDto)
	queryParams.Add("q", q)

	return queryParams.Encode()

}
