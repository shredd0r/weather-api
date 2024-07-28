package http

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

var (
	ErrInvalidLocale     = errors.New("invalid locale")
	ErrInvalidMetric     = errors.New("invalid metric")
	ErrInvalidCoords     = errors.New("invalid coords")
	ErrCountRequestIsOut = errors.New("count of request to api out")
	ErrStatusCodeNot200  = errors.New("api returned code != 200 OK")
)

type Client struct {
	BaseURL    string
	apiKey     string
	log        *logrus.Logger
	httpClient *http.Client
}

func NewHttpClient(log *logrus.Logger) *http.Client {
	return &http.Client{
		Transport: &LoggingTransport{
			log:       log,
			Transport: http.DefaultTransport,
		},
	}
}

func NewClient(baseUrlString string, log *logrus.Logger, httpClient *http.Client, apiKey string) *Client {
	return &Client{
		BaseURL:    baseUrlString,
		apiKey:     apiKey,
		log:        log,
		httpClient: httpClient,
	}
}

func ResponseBodyDecoder[T any](r io.ReadCloser) T {
	defer r.Close()
	var responseBody T

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&responseBody); err != nil {
		panic("Error when decoding response body")
	}

	return responseBody
}

func HttpGetAndGetResponse[T any](httpClient *http.Client, log *logrus.Logger, request *http.Request) (*T, error) {
	response, err := httpClient.Do(request)

	defer response.Body.Close()

	if err != nil {
		log.Infof("error from request: %s", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, ErrStatusCodeNot200
	}

	responseBody := ResponseBodyDecoder[T](response.Body)
	return &responseBody, nil
}

func GetHttpRequestBy(urlForRequest string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, urlForRequest, nil)
	return req
}
