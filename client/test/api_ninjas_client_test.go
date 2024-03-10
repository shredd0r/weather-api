package test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"weather-api/client"
	"weather-api/config"
	"weather-api/dto"
	"weather-api/logger"
)

func TestApiNinjasGetGeocodingInfo(t *testing.T) {
	c := initApiNinjasTest(t)
	response, err := c.GetGeocoding(getApiNinjasGetGeocodingRequest())

	assert.Nil(t, err)
	assert.NotNil(t, response)

}

func getApiNinjasGetGeocodingRequest() dto.ApiNinjasGeocodingRequestDto {
	return dto.ApiNinjasGeocodingRequestDto{
		City: "Kharkiv",
	}
}

func initApiNinjasTest(t *testing.T) *client.ApiNinjasClient {
	cfg := config.ParseEnv()
	if cfg.ApiNinjasApiKey == "" {
		t.Skip("skip test because api key not set")
	}
	log := logger.NewLogger(cfg.Logger)

	return client.NewApiNinjasClient(
		log,
		client.NewHttpClient(log),
		cfg.ApiNinjasApiKey)
}
