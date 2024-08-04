package test

import (
	"github.com/shredd0r/weather-api/client/http"
	"github.com/shredd0r/weather-api/config"
	"github.com/shredd0r/weather-api/dto"
	"github.com/shredd0r/weather-api/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApiNinjasGetGeocodingInfo(t *testing.T) {
	c := initApiNinjasTest(t)
	response, err := c.GetGeocoding(getApiNinjasGetGeocodingRequest())

	assert.Nil(t, err)
	assert.NotNil(t, response)

}

func TestApiNinjasGetReverseGeocodingInfo(t *testing.T) {
	c := initApiNinjasTest(t)
	response, err := c.GetReversGeocoding(getApiNinjasGetReverseGeocodingRequest())

	assert.Nil(t, err)
	assert.NotNil(t, response)

}

func getApiNinjasGetGeocodingRequest() dto.ApiNinjasGeocodingRequestDto {
	return dto.ApiNinjasGeocodingRequestDto{
		City: "Kharkiv",
	}
}

func getApiNinjasGetReverseGeocodingRequest() dto.ApiNinjasReverseGeocodingRequestDto {
	return dto.ApiNinjasReverseGeocodingRequestDto{
		Latitude:  44.5888,
		Longitude: 33.5224,
	}
}

func initApiNinjasTest(t *testing.T) *http.ApiNinjasClient {
	cfg := config.ParseEnv()
	if cfg.ApiKeys.ApiNinjasApiKey == "" {
		t.Skip("skip test because api key not set")
	}
	log := log.NewLogger(cfg.Logger)

	return http.NewApiNinjasClient(
		log,
		http.NewHttpClient(log),
		cfg.ApiKeys.ApiNinjasApiKey)
}
