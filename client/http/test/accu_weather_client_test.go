package test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"weather-api/client/http"
	"weather-api/config"
	"weather-api/dto"
	"weather-api/logger"
)

func TestAccuCurrentWeatherInfo(t *testing.T) {
	c := initAccuWeatherTest(t)
	response, err := c.GetCurrentWeatherInfo(getAccuWeatherRequestDto())

	assert.Nil(t, err)
	assert.NotNil(t, response)

}

func TestAccuHouryWeatherInfo(t *testing.T) {
	c := initAccuWeatherTest(t)
	response, err := c.GetHourlyWeatherInfo(getAccuWeatherRequestDto())

	assert.Nil(t, err)
	assert.NotNil(t, response)
}

func TestAccuDailyWeatherInfo(t *testing.T) {
	c := initAccuWeatherTest(t)
	response, err := c.GetDailyWeatherInfo(getAccuWeatherRequestDto())

	assert.Nil(t, err)
	assert.NotNil(t, response)
}

func TestAccuGeoPositionInfo(t *testing.T) {
	c := initAccuWeatherTest(t)
	response, err := c.GetGeoPositionSearch(getGeoPositionRequest())

	assert.Nil(t, err)
	assert.NotNil(t, response)
}

func getAccuWeatherRequestDto() dto.AccuWeatherRequestDto {
	return dto.AccuWeatherRequestDto{
		AccuWeatherBaseRequestDto: getAccuWeatherBaseRequestDto(),
		LocationKey:               "1216600",
	}
}

func getGeoPositionRequest() dto.AccuWeatherGeoPositionRequestDto {
	return dto.AccuWeatherGeoPositionRequestDto{
		AccuWeatherBaseRequestDto: getAccuWeatherBaseRequestDto(),
		Latitude:                  50.4536,
		Longitude:                 30.5164,
	}
}

func getAccuWeatherBaseRequestDto() dto.AccuWeatherBaseRequestDto {
	return dto.AccuWeatherBaseRequestDto{
		Language: "uk",
		Metric:   true,
		Details:  true,
	}
}

func initAccuWeatherTest(t *testing.T) *http.AccuWeatherClient {
	cfg := config.ParseEnv()
	if cfg.AccuWeatherApiKey == "" {
		t.Skip("skip test because api key not set")
	}
	log := logger.NewLogger(cfg.Logger)

	return http.NewAccuWeatherClient(
		log,
		http.NewHttpClient(log),
		cfg.AccuWeatherApiKey)
}
