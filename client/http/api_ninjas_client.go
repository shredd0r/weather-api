package http

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"weather-api/dto"
	"weather-api/util"
)

type ApiNinjasClient struct {
	*Client
}

var headerKey = "X-Api-Key"

func NewApiNinjasClient(log *logrus.Logger, httpClient *http.Client, apiKey string) *ApiNinjasClient {
	return &ApiNinjasClient{
		NewClient("https://api.api-ninjas.com/", log, httpClient, apiKey),
	}
}

func (c *ApiNinjasClient) GetGeocoding(request dto.ApiNinjasGeocodingRequestDto) (*dto.ApiNinjasGeocodingResponseDto, error) {
	urlForRequest := c.BaseURL + "v1/geocoding?" + c.getQueryParams(request)
	req := GetHttpRequestBy(urlForRequest)
	req.Header.Set(headerKey, c.apiKey)

	response, err := HttpGetAndGetResponse[[]*dto.ApiNinjasGeocodingResponseDto](c.httpClient, c.log, req)

	return GetFirstElemFromResponse[dto.ApiNinjasGeocodingResponseDto](response, err)
}

func (c *ApiNinjasClient) GetReversGeocoding(request dto.ApiNinjasReverseGeocodingRequestDto) (*dto.ApiNinjasReverseGeocodingResponseDto, error) {
	urlForRequest := c.BaseURL + "v1/reversegeocoding?" + c.getQueryParamsForReverse(request)
	req := GetHttpRequestBy(urlForRequest)
	req.Header.Set(headerKey, c.apiKey)

	response, err := HttpGetAndGetResponse[[]*dto.ApiNinjasReverseGeocodingResponseDto](c.httpClient, c.log, req)

	return GetFirstElemFromResponse[dto.ApiNinjasReverseGeocodingResponseDto](response, err)

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
