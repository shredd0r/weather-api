package client

import (
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL    *url.URL
	logger     *log.Logger
	httpClient *http.Client
}

func NewClient(baseUrlString string, logger *log.Logger, httpClient *http.Client) *Client {
	baseUrl, _ := url.Parse(baseUrlString)

	return &Client{
		BaseURL:    baseUrl,
		logger:     logger,
		httpClient: httpClient,
	}
}
