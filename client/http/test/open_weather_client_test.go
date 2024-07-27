package test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"weather-api/client/http"
	"weather-api/config"
	"weather-api/dto"
	"weather-api/log"
)

func TestCurrentWeatherInfo(t *testing.T) {
	c := initOpenWeatherTest(t)
	response, err := c.GetCurrentWeatherInfo(getOpenWeatherRequestDto())

	assert.Nil(t, err)
	assert.NotNil(t, response)

}

func TestForecastInfo(t *testing.T) {
	c := initOpenWeatherTest(t)
	response, err := c.GetForecastWeatherInfo(getOpenWeatherForecastRequestDto())

	assert.Nil(t, err)
	assert.NotNil(t, response)
}

func getOpenWeatherRequestDto() dto.OpenWeatherRequestDto {
	return dto.OpenWeatherRequestDto{
		Latitude:  50.4536,
		Longitude: 30.5164,
		Units:     "metric",
		Language:  "uk"}
}

func getOpenWeatherForecastRequestDto() dto.OpenWeatherForecastRequestDto {
	return dto.OpenWeatherForecastRequestDto{
		OpenWeatherRequestDto: getOpenWeatherRequestDto(),
		Count:                 10,
	}
}

func initOpenWeatherTest(t *testing.T) *http.OpenWeatherClient {
	cfg := config.ParseEnv()
	if cfg.ApiKeys.OpenWeatherApiKey == "" {
		t.Skip("skip test because api key not set")
	}
	log := log.NewLogger(cfg.Logger)

	return http.NewOpenWeatherClient(
		log,
		http.NewHttpClient(log),
		cfg.ApiKeys.OpenWeatherApiKey)
}
