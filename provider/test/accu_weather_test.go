package test

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"weather-api/client/http"
	mock_http "weather-api/client/http/mock"
	"weather-api/config"
	"weather-api/dto"
	"weather-api/logger"
	"weather-api/provider"
)

type callMethod[T any] func(requestDto dto.WeatherRequestDto) (T, error)

func TestCurrentWeather(t *testing.T) {
	p := getHttpAccuWeatherProvider(t)

	resp, err := p.CurrentWeatherInfo(getRequest())

	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestHourlyWeather(t *testing.T) {
	p := getHttpAccuWeatherProvider(t)

	resp, err := p.HourlyWeatherInfo(getRequest())

	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestDailyWeather(t *testing.T) {
	p := getHttpAccuWeatherProvider(t)

	resp, err := p.DailyWeatherInfo(getRequest())

	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestErrorFromClient(t *testing.T) {
	expectedErr := http.StatusCodeNot200
	req := getRequest()
	testCases := []struct {
		name         string
		testedMethod func(p *provider.HttpAccuWeatherProvider, c *mock_http.MockAccuWeatherInterface) (interface{}, error)
	}{
		{
			name: "TestErrorFromClientCurrentWeather",
			testedMethod: func(p *provider.HttpAccuWeatherProvider, c *mock_http.MockAccuWeatherInterface) (interface{}, error) {
				c.EXPECT().GetCurrentWeatherInfo(gomock.Any()).Return(nil, expectedErr)
				return p.CurrentWeatherInfo(req)
			},
		},
		{
			name: "TestErrorFromClientHourlyWeather",
			testedMethod: func(p *provider.HttpAccuWeatherProvider, c *mock_http.MockAccuWeatherInterface) (interface{}, error) {
				c.EXPECT().GetHourlyWeatherInfo(gomock.Any()).Return(nil, expectedErr)
				return p.HourlyWeatherInfo(req)
			},
		},
		{
			name: "TestErrorFromClientDailyWeather",
			testedMethod: func(p *provider.HttpAccuWeatherProvider, c *mock_http.MockAccuWeatherInterface) (interface{}, error) {
				c.EXPECT().GetDailyWeatherInfo(gomock.Any()).Return(nil, expectedErr)
				return p.DailyWeatherInfo(req)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			p, c := getProviderAndMock(t)

			resp, err := testCase.testedMethod(p, c)

			assert.ErrorIs(t, err, expectedErr)
			assert.Nil(t, resp)
		})
	}
}

func TestNilInResponse(t *testing.T) {
	unitSlice := []dto.Unit{dto.UnitImperial, dto.UnitMetric}
	valueWithNil := dto.AccuWeatherValueInfoDto{Value: nil}
	indicationWithNil := dto.AccuWeatherIndicationInfoDto{Metric: valueWithNil, Imperial: valueWithNil}
	windWithNil := dto.AccuWeatherWindInfoDto{Speed: valueWithNil}
	dayWithNil := dto.AccuWeatherDayInfoDto{
		PrecipitationProbability: nil,
		TotalLiquid:              valueWithNil,
		Wind:                     windWithNil,
	}

	testCases := []struct {
		name         string
		testedMethod func(*provider.HttpAccuWeatherProvider, *mock_http.MockAccuWeatherInterface, dto.Unit) (interface{}, error)
	}{
		{
			name: "TestNilExistInCurrentWeatherResponse",
			testedMethod: func(p *provider.HttpAccuWeatherProvider, c *mock_http.MockAccuWeatherInterface, unit dto.Unit) (interface{}, error) {
				c.EXPECT().GetCurrentWeatherInfo(gomock.Any()).Return(&dto.AccuWeatherCurrentResponseDto{
					WeatherIcon:         nil,
					Temperature:         indicationWithNil,
					RealFeelTemperature: indicationWithNil,
					Visibility:          indicationWithNil,
					UVIndex:             nil,
				}, nil)
				return p.CurrentWeatherInfo(getRequestBy(unit))
			},
		},
		{
			name: "TestNilExistInHourlyWeatherResponse",
			testedMethod: func(p *provider.HttpAccuWeatherProvider, c *mock_http.MockAccuWeatherInterface, unit dto.Unit) (interface{}, error) {
				c.EXPECT().GetHourlyWeatherInfo(gomock.Any()).Return(&[]dto.AccuWeatherHourlyResponseDto{
					{
						WeatherIcon:              nil,
						Temperature:              valueWithNil,
						RealFeelTemperature:      valueWithNil,
						Wind:                     windWithNil,
						UVIndex:                  nil,
						PrecipitationProbability: nil,
						RainProbability:          nil,
						SnowProbability:          nil,
						IceProbability:           nil,
						TotalLiquid:              valueWithNil,
						Rain:                     valueWithNil,
						Snow:                     valueWithNil,
						Ice:                      valueWithNil,
					},
				}, nil)
				return p.HourlyWeatherInfo(getRequestBy(unit))
			},
		},
		{
			name: "TestNilExistInDailyWeatherResponse",
			testedMethod: func(p *provider.HttpAccuWeatherProvider, c *mock_http.MockAccuWeatherInterface, unit dto.Unit) (interface{}, error) {
				c.EXPECT().GetDailyWeatherInfo(gomock.Any()).Return(&dto.AccuWeatherDailyResponseDto{
					DailyForecasts: []dto.AccuWeatherDailyForecastDto{
						{
							Temperature:  dto.AccuWeatherTemperatureDto{Minimum: valueWithNil, Maximum: valueWithNil},
							AirAndPollen: []dto.AccuWeatherCategoryInfoDto{{Value: nil}},
							Day:          dayWithNil,
							Night:        dayWithNil,
						},
					},
				}, nil)
				return p.DailyWeatherInfo(getRequestBy(unit))
			},
		},
	}

	for _, unit := range unitSlice {
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				p, c := getProviderAndMock(t)

				resp, err := testCase.testedMethod(p, c, unit)

				assert.Nil(t, err)
				assert.NotNil(t, resp)
			})
		}
	}

}

func getHttpAccuWeatherProvider(t *testing.T) *provider.HttpAccuWeatherProvider {
	cfg := config.ParseEnv()
	if cfg.AccuWeatherApiKey == "" {
		t.Skip("skip test because accu weather api key not set in env")
	}

	log := logger.NewLogger(cfg.Logger)
	httpClient := http.NewHttpClient(log)
	client := http.NewAccuWeatherClient(log, httpClient, cfg.AccuWeatherApiKey)
	return provider.NewHttpAccuWeatherProvider(client, log)
}

func getProviderAndMock(t *testing.T) (*provider.HttpAccuWeatherProvider, *mock_http.MockAccuWeatherInterface) {
	ctrl := gomock.NewController(t)
	log := logger.StandardLogger()
	client := mock_http.NewMockAccuWeatherInterface(ctrl)
	p := provider.NewHttpAccuWeatherProvider(client, log)

	return p, client
}

func getRequest() dto.WeatherRequestDto {
	return getRequestBy(dto.UnitMetric)
}

func getRequestBy(unit dto.Unit) dto.WeatherRequestDto {
	return dto.WeatherRequestDto{
		Unit:     unit,
		Language: dto.LanguageUkrainian,
		LocalityDto: dto.LocalityDto{
			AccuWeatherLocationKey: "1216600",
		},
	}
}
