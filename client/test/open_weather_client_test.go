package test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"weather-api/client"
	"weather-api/config"
	"weather-api/dto"
	"weather-api/logger"
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

func initOpenWeatherTest(t *testing.T) *client.OpenWeatherClient {
	cfg := config.ParseEnv()
	if cfg.OpenWeather.ApiKey == "" {
		t.Skip("skip test because api key not set")
	}
	log := logger.NewLogger(cfg.Logger)

	return client.NewOpenWeatherClient(
		log,
		client.NewHttpClient(log),
		&cfg.OpenWeather)
}
