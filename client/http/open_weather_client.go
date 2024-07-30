package http

import (
	"net/http"
	"net/url"
	"strconv"
	"weather-api/dto"
	"weather-api/log"
	"weather-api/util"
)

var (
	errWrongLatitude  = "wrong latitude"
	errWrongLongitude = "wrong longitude"
)

type OpenWeatherInterface interface {
	GetCurrentWeatherInfo(request dto.OpenWeatherRequestDto) (*dto.OpenWeatherCurrentWeatherResponseDto, error)
	GetForecastWeatherInfo(request dto.OpenWeatherForecastRequestDto) (*dto.OpenWeatherHourlyWeatherResponseDto, error)
}

type OpenWeatherClient struct {
	client *Client
}

func NewOpenWeatherClient(log log.Logger, httpClient *http.Client, apiKey string) *OpenWeatherClient {
	return &OpenWeatherClient{
		NewClient("https://api.openweathermap.org/", log, httpClient, apiKey),
	}
}

func (c *OpenWeatherClient) GetCurrentWeatherInfo(request dto.OpenWeatherRequestDto) (*dto.OpenWeatherCurrentWeatherResponseDto, error) {
	urlForRequest := c.client.BaseURL + "data/2.5/weather?" + c.getQueryParamByRequest(request).Encode()
	return HttpGetAndGetResponse[dto.OpenWeatherCurrentWeatherResponseDto](
		c.client.httpClient,
		c.client.Log,
		GetHttpRequestBy(urlForRequest),
		c.parseErrorInResponse)
}

func (c *OpenWeatherClient) GetForecastWeatherInfo(request dto.OpenWeatherForecastRequestDto) (*dto.OpenWeatherHourlyWeatherResponseDto, error) {
	urlForRequest := c.client.BaseURL + "data/2.5/forecast?" + c.getQueryParamsForForecast(request)
	return HttpGetAndGetResponse[dto.OpenWeatherHourlyWeatherResponseDto](
		c.client.httpClient,
		c.client.Log,
		GetHttpRequestBy(urlForRequest),
		c.parseErrorInResponse)
}

func (c *OpenWeatherClient) getQueryParamByRequest(request dto.OpenWeatherRequestDto) url.Values {
	queryParams := url.Values{
		"appid": {c.client.apiKey},
		"lat":   {util.Float64ToString(request.Latitude)},
		"lon":   {util.Float64ToString(request.Longitude)},
		"units": {string(request.Units)},
		"lang":  {request.Language},
	}

	return queryParams
}

func (c *OpenWeatherClient) getQueryParamsForForecast(request dto.OpenWeatherForecastRequestDto) string {
	queryParams := c.getQueryParamByRequest(request.OpenWeatherRequestDto)
	if request.Count != nil {
		queryParams.Add("cnt", strconv.Itoa(*request.Count))
	}

	return queryParams.Encode()
}

func (c *OpenWeatherClient) parseErrorInResponse(resp *http.Response) error {
	openWeatherError, err := ResponseBodyDecoder[dto.OpenWeatherError](resp.Body)

	if err != nil {
		return err
	}

	if openWeatherError.Cod == http.StatusUnauthorized {
		return ErrCountRequestIsOut
	}

	switch openWeatherError.Message {
	case errWrongLatitude, errWrongLongitude:
		{
			return ErrInvalidCoords
		}
	}

	return nil
}
