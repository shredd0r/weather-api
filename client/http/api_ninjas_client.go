package http

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"weather-api/dto"
	"weather-api/util"
)

type ApiNinjasInterface interface {
	GetGeocoding(request dto.ApiNinjasGeocodingRequestDto) (*[]*dto.ApiNinjasGeocodingResponseDto, error)
	GetReversGeocoding(request dto.ApiNinjasReverseGeocodingRequestDto) (*[]*dto.ApiNinjasReverseGeocodingResponseDto, error)
}

type ApiNinjasClient struct {
	client *Client
}

var headerKey = "X-Api-Key"

func NewApiNinjasClient(log *logrus.Logger, httpClient *http.Client, apiKey string) *ApiNinjasClient {
	return &ApiNinjasClient{
		NewClient("https://api.api-ninjas.com/", log, httpClient, apiKey),
	}
}

func (c *ApiNinjasClient) GetGeocoding(request dto.ApiNinjasGeocodingRequestDto) (*[]*dto.ApiNinjasGeocodingResponseDto, error) {
	urlForRequest := c.client.BaseURL + "v1/geocoding?" + c.getQueryParams(request)
	req := GetHttpRequestBy(urlForRequest)
	req.Header.Set(headerKey, c.client.apiKey)

	response, err := HttpGetAndGetResponse[[]*dto.ApiNinjasGeocodingResponseDto](c.client.httpClient, c.client.log, req)

	return response, err
}

func (c *ApiNinjasClient) GetReversGeocoding(request dto.ApiNinjasReverseGeocodingRequestDto) (*[]*dto.ApiNinjasReverseGeocodingResponseDto, error) {
	urlForRequest := c.client.BaseURL + "v1/reversegeocoding?" + c.getQueryParamsForReverse(request)
	req := GetHttpRequestBy(urlForRequest)
	req.Header.Set(headerKey, c.client.apiKey)

	response, err := HttpGetAndGetResponse[[]*dto.ApiNinjasReverseGeocodingResponseDto](c.client.httpClient, c.client.log, req)

	return response, err

}

func (c *ApiNinjasClient) getQueryParams(request dto.ApiNinjasGeocodingRequestDto) string {
	return url.Values{
		"city":    {request.City},
		"state":   {request.State},
		"country": {request.Country},
	}.Encode()
}

func (c *ApiNinjasClient) getQueryParamsForReverse(request dto.ApiNinjasReverseGeocodingRequestDto) string {
	return url.Values{
		"lat": {util.Float64ToString(request.Latitude)},
		"lon": {util.Float64ToString(request.Longitude)},
	}.Encode()
}
