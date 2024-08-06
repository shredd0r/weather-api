package http

import (
	"github.com/shredd0r/weather-api/dto"
	"github.com/shredd0r/weather-api/log"
	"github.com/shredd0r/weather-api/util"
	"net/http"
	"net/url"
)

var (
	errMissingAPIKey = "Missing API Key."
)

type ApiNinjasInterface interface {
	GetGeocoding(request dto.ApiNinjasGeocodingRequestDto) (*[]*dto.ApiNinjasGeocodingResponseDto, error)
	GetReversGeocoding(request dto.ApiNinjasReverseGeocodingRequestDto) (*[]*dto.ApiNinjasReverseGeocodingResponseDto, error)
}

type ApiNinjasClient struct {
	client *Client
}

var headerKey = "X-Api-Key"

func NewApiNinjasClient(log log.Logger, httpClient *http.Client, apiKey string) *ApiNinjasClient {
	return &ApiNinjasClient{
		NewClient("https://api.api-ninjas.com/", log, httpClient, apiKey),
	}
}

func (c *ApiNinjasClient) GetGeocoding(request dto.ApiNinjasGeocodingRequestDto) (*[]*dto.ApiNinjasGeocodingResponseDto, error) {
	urlForRequest := c.client.BaseURL + "v1/geocoding?" + c.getQueryParams(request)
	req := GetHttpRequestBy(urlForRequest)
	req.Header.Set(headerKey, c.client.apiKey)

	response, err := HttpGetAndGetResponse[[]*dto.ApiNinjasGeocodingResponseDto](
		c.client.httpClient,
		c.client.Log,
		req,
		c.parseErrorInResponse)

	return response, err
}

func (c *ApiNinjasClient) GetReversGeocoding(request dto.ApiNinjasReverseGeocodingRequestDto) (*[]*dto.ApiNinjasReverseGeocodingResponseDto, error) {
	urlForRequest := c.client.BaseURL + "v1/reversegeocoding?" + c.getQueryParamsForReverse(request)
	req := GetHttpRequestBy(urlForRequest)
	req.Header.Set(headerKey, c.client.apiKey)

	response, err := HttpGetAndGetResponse[[]*dto.ApiNinjasReverseGeocodingResponseDto](
		c.client.httpClient,
		c.client.Log,
		req,
		c.parseErrorInResponse)

	return response, err

}

func (c *ApiNinjasClient) getQueryParams(request dto.ApiNinjasGeocodingRequestDto) string {
	values := url.Values{}

	if request.City != nil {
		values.Add("city", *request.City)
	}

	if request.State != nil {
		values.Add("state", *request.State)
	}

	if request.Country != nil {
		values.Add("country", *request.Country)
	}

	return values.Encode()
}

func (c *ApiNinjasClient) getQueryParamsForReverse(request dto.ApiNinjasReverseGeocodingRequestDto) string {
	return url.Values{
		"lat": {util.Float64ToString(request.Latitude)},
		"lon": {util.Float64ToString(request.Longitude)},
	}.Encode()
}

func (c *ApiNinjasClient) parseErrorInResponse(resp *http.Response) error {
	_, err := ResponseBodyDecoder[dto.ApiNinjasError](resp.Body)

	if err != nil {
		return ErrCountRequestIsOut
	}

	return nil
}
