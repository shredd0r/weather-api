package http

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"weather-api/dto"
	"weather-api/log"
)

var (
	errMessageAuthorizationFailed = "Api Authorization failed"
)

type AccuWeatherInterface interface {
	GetCurrentWeatherInfo(request dto.AccuWeatherRequestDto) (*dto.AccuWeatherCurrentResponseDto, error)
	GetHourlyWeatherInfo(request dto.AccuWeatherRequestDto) (*[]*dto.AccuWeatherHourlyResponseDto, error)
	GetDailyWeatherInfo(request dto.AccuWeatherRequestDto) (*dto.AccuWeatherDailyResponseDto, error)
	GetGeoPositionSearch(request dto.AccuWeatherGeoPositionRequestDto) (*dto.AccuWeatherGeoPositionResponseDto, error)
}

type AccuWeatherClient struct {
	client *Client
}

func NewAccuWeatherClient(log log.Logger, httpClient *http.Client, apiKey string) *AccuWeatherClient {
	return &AccuWeatherClient{
		NewClient("http://dataservice.accuweather.com/", log, httpClient, apiKey),
	}
}

func (c *AccuWeatherClient) GetCurrentWeatherInfo(request dto.AccuWeatherRequestDto) (*dto.AccuWeatherCurrentResponseDto, error) {
	var urlForRequest = c.client.BaseURL + fmt.Sprintf("currentconditions/v1/%s?", request.LocationKey) + c.getQueryParamsForRequest(request)
	response, err := HttpGetAndGetResponse[[]*dto.AccuWeatherCurrentResponseDto](
		c.client.httpClient,
		c.client.Log,
		GetHttpRequestBy(urlForRequest),
		c.parseErrorInResponse)

	if err == nil {
		return (*response)[0], err
	}

	c.client.Log.Error("error when get current weather info")

	return nil, err
}

func (c *AccuWeatherClient) GetHourlyWeatherInfo(request dto.AccuWeatherRequestDto) (*[]*dto.AccuWeatherHourlyResponseDto, error) {
	var urlForRequest = c.client.BaseURL + fmt.Sprintf("forecasts/v1/hourly/12hour/%s?", request.LocationKey) + c.getQueryParamsForRequest(request)
	resp, err := HttpGetAndGetResponse[[]*dto.AccuWeatherHourlyResponseDto](
		c.client.httpClient,
		c.client.Log,
		GetHttpRequestBy(urlForRequest),
		c.parseErrorInResponse)

	if err != nil {
		c.client.Log.Error("error when get hourly weather info")
		return nil, err
	}

	return resp, nil
}

func (c *AccuWeatherClient) GetDailyWeatherInfo(request dto.AccuWeatherRequestDto) (*dto.AccuWeatherDailyResponseDto, error) {
	var urlForRequest = c.client.BaseURL + fmt.Sprintf("forecasts/v1/daily/5day/%s?", request.LocationKey) + c.getQueryParamsForRequest(request)
	resp, err := HttpGetAndGetResponse[dto.AccuWeatherDailyResponseDto](
		c.client.httpClient,
		c.client.Log,
		GetHttpRequestBy(urlForRequest),
		c.parseErrorInResponse)

	if err != nil {
		c.client.Log.Error("error when get daily weather info")
		return nil, err
	}

	return resp, nil

}

func (c *AccuWeatherClient) GetGeoPositionSearch(request dto.AccuWeatherGeoPositionRequestDto) (*dto.AccuWeatherGeoPositionResponseDto, error) {
	var urlForRequest = c.client.BaseURL + "locations/v1/cities/geoposition/search?" + c.getQueryParamsForGeoPosition(request)
	return HttpGetAndGetResponse[dto.AccuWeatherGeoPositionResponseDto](
		c.client.httpClient,
		c.client.Log,
		GetHttpRequestBy(urlForRequest),
		c.parseErrorInResponse)
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

func (c *AccuWeatherClient) parseErrorInResponse(resp *http.Response) error {
	accuWeatherError, err := ResponseBodyDecoder[dto.AccuWeatherError](resp.Body)

	if err != nil {
		return err
	}

	switch accuWeatherError.Message {
	case errMessageAuthorizationFailed:
		{
			return ErrCountRequestIsOut
		}
	}

	return nil
}
