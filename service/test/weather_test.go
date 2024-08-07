package test

import (
	"context"
	"fmt"
	"github.com/shredd0r/weather-api/config"
	"github.com/shredd0r/weather-api/dto"
	"github.com/shredd0r/weather-api/log"
	mock_provider "github.com/shredd0r/weather-api/provider/mock"
	"github.com/shredd0r/weather-api/service"
	mock_service "github.com/shredd0r/weather-api/service/mock"
	"github.com/shredd0r/weather-api/storage"
	mock_storage "github.com/shredd0r/weather-api/storage/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

// TODO tests for concrete weather is copy paste with chance methods for get weather, make more general this test
func TestForCurrentWeather(t *testing.T) {
	logger, locationService, weatherProvider, weatherStorage := getMocks(t)
	for _, forecaster := range getForecasterList() {
		ws := service.NewWeatherService(logger, forecaster, locationService, weatherProvider, weatherStorage)
		testedMethod := ws.CurrentWeather
		returnedWeather := &dto.CurrentWeather{EpochTime: time.Now().UnixMilli()}

		t.Run(
			"Current weather store in storage",
			func(t *testing.T) {
				testWeatherInStorage[dto.CurrentWeather](
					t,
					locationService,
					weatherProvider,
					testedMethod,
					weatherStorage.EXPECT().GetCurrentWeather,
					returnedWeather,
				)
			})

		t.Run(
			"Current weather not found in storage",
			func(t *testing.T) {
				testWeatherNotFoundInStorage[dto.CurrentWeather](
					t, locationService,
					testedMethod,
					weatherStorage.EXPECT().GetCurrentWeather,
					weatherProvider.EXPECT().CurrentWeather,
					returnedWeather,
				)
			},
		)

	}
}

func TestForHourlyWeather(t *testing.T) {
	logger, locationService, weatherProvider, weatherStorage := getMocks(t)
	for _, forecaster := range getForecasterList() {
		ws := service.NewWeatherService(logger, forecaster, locationService, weatherProvider, weatherStorage)
		testedMethod := ws.HourlyWeather
		returnedWeather := &[]*dto.HourlyWeather{}

		t.Run(
			"Hourly weather store in storage",
			func(t *testing.T) {
				testWeatherInStorage[[]*dto.HourlyWeather](
					t,
					locationService,
					weatherProvider,
					testedMethod,
					weatherStorage.EXPECT().GetHourlyWeather,
					returnedWeather,
				)
			})

		t.Run(
			"Hourly weather not found in storage",
			func(t *testing.T) {
				testWeatherNotFoundInStorage[[]*dto.HourlyWeather](
					t, locationService,
					testedMethod,
					weatherStorage.EXPECT().GetHourlyWeather,
					weatherProvider.EXPECT().HourlyWeather,
					returnedWeather,
				)
			},
		)

	}
}

func TestForDailyWeather(t *testing.T) {
	logger, locationService, weatherProvider, weatherStorage := getMocks(t)
	for _, forecaster := range getForecasterList() {
		ws := service.NewWeatherService(logger, forecaster, locationService, weatherProvider, weatherStorage)
		testedMethod := ws.DailyWeather
		returnedWeather := &[]*dto.DailyWeather{}

		t.Run(
			"Daily weather store in storage",
			func(t *testing.T) {
				testWeatherInStorage[[]*dto.DailyWeather](
					t,
					locationService,
					weatherProvider,
					testedMethod,
					weatherStorage.EXPECT().GetDailyWeather,
					returnedWeather,
				)
			})

		t.Run(
			"Daily weather not found in storage",
			func(t *testing.T) {
				testWeatherNotFoundInStorage[[]*dto.DailyWeather](
					t, locationService,
					testedMethod,
					weatherStorage.EXPECT().GetDailyWeather,
					weatherProvider.EXPECT().DailyWeather,
					returnedWeather,
				)
			},
		)

	}
}

func testWeatherNotFoundInStorage[T any](
	t *testing.T,
	locationService *mock_service.MockLocationService,
	testedMethodGetWeatherFromService func(ctx context.Context, request *dto.WeatherRequest) (*T, error),
	expectMethodForGetWeatherFromStorage func(ctx, addressHash, forecaster any) *gomock.Call,
	expectMethodForGetWeatherFromProvider func(ctx, request any) *gomock.Call,
	returnedWeatherFromProvider *T,
) {
	for _, forecaster := range getForecasterList() {
		t.Run(fmt.Sprintf("case for %s", forecaster), func(*testing.T) {
			ctx := context.Background()
			req := getRequest()
			locationService.EXPECT().LocationByCoords(gomock.Any(), req.Coords).Return(&dto.LocationInfo{}, nil).Times(1)
			expectMethodForGetWeatherFromStorage(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, storage.ErrNotFound).Times(1)
			expectMethodForGetWeatherFromProvider(gomock.Any(), gomock.Any()).Return(returnedWeatherFromProvider, nil).Times(1)

			resp, err := testedMethodGetWeatherFromService(ctx, req)

			assert.Nil(t, err)
			assert.Equal(t, resp, returnedWeatherFromProvider)
		})
	}
}

func testWeatherInStorage[T any](
	t *testing.T,
	locationService *mock_service.MockLocationService,
	weatherProvider *mock_provider.MockWeatherProvider,
	testedMethodGetWeatherFromService func(ctx context.Context, request *dto.WeatherRequest) (*T, error),
	expectMethodForGetWeather func(ctx, addressHash, forecaster any) *gomock.Call,
	returnedWeatherDto *T) {

	for _, forecaster := range getForecasterList() {
		t.Run(fmt.Sprintf("case for %s", forecaster), func(*testing.T) {
			ctx := context.Background()
			req := getRequest()
			locationService.EXPECT().LocationByCoords(gomock.Any(), req.Coords).Return(&dto.LocationInfo{}, nil).Times(1)
			expectMethodForGetWeather(gomock.Any(), gomock.Any(), gomock.Any()).Return(returnedWeatherDto, nil).Times(1)
			weatherProvider.EXPECT().HourlyWeather(gomock.Any(), gomock.Any()).Times(0)

			resp, err := testedMethodGetWeatherFromService(ctx, req)

			assert.Nil(t, err)
			assert.NotNil(t, resp)
		})
	}
}

func getMocks(t *testing.T) (log.Logger, *mock_service.MockLocationService, *mock_provider.MockWeatherProvider, *mock_storage.MockWeatherStorage) {
	cfg := config.ParseEnv()
	logger := log.NewLogger(cfg.Logger)

	ctrl := gomock.NewController(t)
	ls := mock_service.NewMockLocationService(ctrl)
	wp := mock_provider.NewMockWeatherProvider(ctrl)
	wst := mock_storage.NewMockWeatherStorage(ctrl)

	return logger, ls, wp, wst
}

func getForecasterList() []dto.WeatherForecaster {
	return []dto.WeatherForecaster{
		dto.WeatherForecasterAccuWeather,
		dto.WeatherForecasterOpenWeather,
	}
}

func getRequest() *dto.WeatherRequest {
	return &dto.WeatherRequest{
		Coords: &dto.Coords{
			Latitude:  9.11,
			Longitude: 9.11,
		},
		Locale: "us",
		Unit:   dto.UnitImperial,
	}
}
