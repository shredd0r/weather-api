package client

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Client struct {
	BaseURL    string
	logger     *logrus.Logger
	httpClient *http.Client
}

func NewClient(baseUrlString string, logger *logrus.Logger, httpClient *http.Client) *Client {
	return &Client{
		BaseURL:    baseUrlString,
		logger:     logger,
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
