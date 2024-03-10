package client

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"weather-api/dto"
)

type ApiNinjasClient struct {
	*Client
}

func NewApiNinjasClient(log *logrus.Logger, httpClient *http.Client, apiKey string) *ApiNinjasClient {
	return &ApiNinjasClient{
		NewClient("https://api.api-ninjas.com/", log, httpClient, apiKey),
	}
}

func (c *ApiNinjasClient) GetGeocoding(request dto.ApiNinjasGeocodingRequestDto) (*[]*dto.ApiNinjasGeocodingResponseDto, error) {
	urlForRequest := c.BaseURL + "v1/geocoding?" + c.getQueryParams(request)
	req := GetHttpRequestBy(urlForRequest)
	req.Header.Set("X-Api-Key", c.apiKey)

	return HttpGetAndGetResponse[[]*dto.ApiNinjasGeocodingResponseDto](c.httpClient, c.log, req)
}

func (c *ApiNinjasClient) getQueryParams(request dto.ApiNinjasGeocodingRequestDto) string {
	return url.Values{
		"city":    {request.City},
		"state":   {request.State},
		"country": {request.Country},
	}.Encode()
}
