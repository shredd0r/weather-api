package http

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
	"weather-api/dto"
)

type OpenWeatherClient struct {
	*Client
}

func NewOpenWeatherClient(log *logrus.Logger, httpClient *http.Client, apiKey string) *OpenWeatherClient {
	return &OpenWeatherClient{
		NewClient("https://api.openweathermap.org/", log, httpClient, apiKey),
	}
}

func (c *OpenWeatherClient) GetCurrentWeatherInfo(request dto.OpenWeatherRequestDto) (*dto.OpenWeatherCurrentWeatherResponseDto, error) {
	urlForRequest := c.BaseURL + "data/2.5/weather?" + c.getQueryParamByRequest(request).Encode()
	return HttpGetAndGetResponse[dto.OpenWeatherCurrentWeatherResponseDto](c.httpClient, c.log, GetHttpRequestBy(urlForRequest))
}

func (c *OpenWeatherClient) GetForecastWeatherInfo(request dto.OpenWeatherForecastRequestDto) (*dto.OpenWeatherHourlyWeatherResponseDto, error) {
	urlForRequest := c.BaseURL + "data/2.5/forecast?" + c.getQueryParamsForForecast(request)
	return HttpGetAndGetResponse[dto.OpenWeatherHourlyWeatherResponseDto](c.httpClient, c.log, GetHttpRequestBy(urlForRequest))

}

func (c *OpenWeatherClient) getQueryParamByRequest(request dto.OpenWeatherRequestDto) url.Values {
	queryParams := url.Values{
		"appid": {c.apiKey},
		"lat":   {strconv.FormatFloat(request.Latitude, 'f', -1, 64)},
		"lon":   {strconv.FormatFloat(request.Longitude, 'f', -1, 64)},
		"units": {request.Units},
		"lang":  {request.Language},
	}

	return queryParams
}

func (c *OpenWeatherClient) getQueryParamsForForecast(request dto.OpenWeatherForecastRequestDto) string {
	queryParams := c.getQueryParamByRequest(request.OpenWeatherRequestDto)
	queryParams.Add("cnt", strconv.FormatInt(int64(request.Count), request.Count))

	return queryParams.Encode()
}
