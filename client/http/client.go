package http

import (
	"encoding/json"
	"errors"
	"github.com/shredd0r/weather-api/log"
	"io"
	"net/http"
)

var (
	ErrInvalidCoords     = errors.New("invalid coords")
	ErrCountRequestIsOut = errors.New("count of request to api out")
)

type Client struct {
	BaseURL    string
	apiKey     string
	Log        log.Logger
	httpClient *http.Client
}

func NewHttpClient(log log.Logger) *http.Client {
	return &http.Client{
		Transport: &LoggingTransport{
			log:       log,
			Transport: http.DefaultTransport,
		},
	}
}

func NewClient(baseUrlString string, log log.Logger, httpClient *http.Client, apiKey string) *Client {
	return &Client{
		BaseURL:    baseUrlString,
		apiKey:     apiKey,
		Log:        log,
		httpClient: httpClient,
	}
}

func ResponseBodyDecoder[T any](r io.ReadCloser) (*T, error) {
	var responseBody T

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&responseBody); err != nil {
		return nil, err
	}

	return &responseBody, nil
}

func HttpGetAndGetResponse[T any](httpClient *http.Client, log log.Logger, request *http.Request, parserError func(resp *http.Response) error) (*T, error) {
	response, err := httpClient.Do(request)

	defer response.Body.Close()

	if err != nil {
		log.Infof("error from request: %s", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		err = parserError(response)
		return nil, err
	}

	return ResponseBodyDecoder[T](response.Body)
}

func GetHttpRequestBy(urlForRequest string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, urlForRequest, nil)
	return req
}
