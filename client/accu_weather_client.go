package client

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
	"weather-api/dto"
)

type AccuWeatherClient struct {
	*Client
}

func NewAccuWeatherClient(logger *logrus.Logger, httpClient *http.Client) *AccuWeatherClient {
	return &AccuWeatherClient{
		NewClient("http://dataservice.accuweather.com/", logger, httpClient),
	}
}

func (c *AccuWeatherClient) GetCurrentWeatherInfo(request dto.AccuWeatherRequestDto) (*dto.AccuWeatherCurrentResponseDto, error) {
	var urlForRequest = c.BaseURL + fmt.Sprintf("currentconditions/v1/%s?", request.LocationKey) + getQueryParamsForRequest(request)
	response, err := c.httpClient.Get(urlForRequest)

	defer response.Body.Close()

	if err != nil {
		c.logger.Infof("error from request: %s", err)
		return nil, err
	}

	responseBody := ResponseBodyDecoder[[]dto.AccuWeatherCurrentResponseDto](response.Body)

	return &responseBody[0], nil
}

func (c *AccuWeatherClient) GetHourlyWeatherInfo(request dto.AccuWeatherRequestDto) ([]*dto.AccuWeatherHourlyResponseDto, error) {
	var urlForRequest = c.BaseURL + fmt.Sprintf("forecasts/v1/hourly/12hour/%s?", request.LocationKey) + getQueryParamsForRequest(request)
	response, err := c.httpClient.Get(urlForRequest)

	defer response.Body.Close()

	if err != nil {
		c.logger.Infof("error from request: %s", err)
		return nil, err
	}

	responseBody := ResponseBodyDecoder[[]*dto.AccuWeatherHourlyResponseDto](response.Body)

	return responseBody, nil
}

func (c *AccuWeatherClient) GetDailyWeatherInfo(request dto.AccuWeatherRequestDto) (*dto.AccuWeatherDailyResponseDto, error) {
	var urlForRequest = c.BaseURL + fmt.Sprintf("forecasts/v1/daily/5day/%s?", request.LocationKey) + getQueryParamsForRequest(request)
	response, err := c.httpClient.Get(urlForRequest)

	defer response.Body.Close()

	if err != nil {
		c.logger.Infof("error from request: %s", err)
		return nil, err
	}

	responseBody := ResponseBodyDecoder[*dto.AccuWeatherDailyResponseDto](response.Body)

	return responseBody, nil
}

func (c *AccuWeatherClient) GetGeoPositionSearch(request dto.AccuWeatherGeoPositionRequestDto) (*dto.AccuWeatherGeoPositionResponseDto, error) {
	var urlForRequest = c.BaseURL + "locations/v1/cities/geoposition/search?" + getQueryParamsForGeoPosition(request)
	response, err := c.httpClient.Get(urlForRequest)

	defer response.Body.Close()

	if err != nil {
		c.logger.Infof("error from request: %s", err)
		return nil, err
	}

	responseBody := ResponseBodyDecoder[dto.AccuWeatherGeoPositionResponseDto](response.Body)

	return &responseBody, nil
}

func getQueryParamsForBase(request dto.AccuWeatherBaseRequestDto) url.Values {
	queryParams := url.Values{
		"apikey":   {request.AppKey},
		"language": {request.Language},
		"details":  {strconv.FormatBool(request.Details)},
		"metric":   {strconv.FormatBool(request.Metric)},
	}
	return queryParams
}

func getQueryParamsForRequest(request dto.AccuWeatherRequestDto) string {
	queryParams := getQueryParamsForBase(request.AccuWeatherBaseRequestDto)
	queryParams.Add("locationKey", request.LocationKey)

	return queryParams.Encode()
}

func getQueryParamsForGeoPosition(request dto.AccuWeatherGeoPositionRequestDto) string {
	q := fmt.Sprintf("%.2f,%.2f", request.Latitude, request.Longitude)
	queryParams := getQueryParamsForBase(request.AccuWeatherBaseRequestDto)
	queryParams.Add("q", q)

	return queryParams.Encode()

}
