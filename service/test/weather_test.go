package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/shredd0r/weather-api/client/http"
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
)

func TestForCurrentWeather(t *testing.T) {
	logger, locationService, weatherProvider, weatherStorage := getMocks(t)
	for _, forecaster := range getForecasterList() {
		ws := service.NewWeatherService(logger, forecaster, locationService, weatherProvider, weatherStorage)

		inputDatas := getInputDatas(ws, weatherStorage, weatherProvider)

		for _, data := range *inputDatas {
			t.Run(
				fmt.Sprintf("%s weather store in storage for forecaster: %s", data.weatherType, forecaster),
				func(t *testing.T) {
					testWeatherInStorage(
						t,
						locationService,
						data.testedMethod,
						data.expectMethodForGetWeatherFromStorage,
						data.expectMethodForGetWeatherFromProvider,
						data.expectMethodForSaveWeather,
						data.expectMethodForSaveUpdatedTimeWeather,
						forecaster,
						data.returnedWeather,
					)
				})

			t.Run(
				fmt.Sprintf("%s weather not found in storage for forecaster: %s", data.weatherType, forecaster),
				func(t *testing.T) {
					testWeatherNotFoundInStorage(
						t,
						locationService,
						data.testedMethod,
						data.expectMethodForGetWeatherFromStorage,
						data.expectMethodForGetWeatherFromProvider,
						data.expectMethodForSaveWeather,
						data.expectMethodForSaveUpdatedTimeWeather,
						forecaster,
						data.returnedWeather,
					)
				},
			)

			t.Run(
				fmt.Sprintf("%s weather provider from forecaster return error for forecaster: %s", data.weatherType, forecaster),
				func(t *testing.T) {
					testWeatherProviderForecasterReturnError(
						t,
						locationService,
						data.testedMethod,
						data.expectMethodForGetWeatherFromStorage,
						data.expectMethodForGetWeatherFromProvider,
						data.expectMethodForSaveWeather,
						data.expectMethodForSaveUpdatedTimeWeather,
						forecaster,
					)
				})

			t.Run(
				fmt.Sprintf("%s invalid coords in request for forecaster: %s", data.weatherType, forecaster),
				func(t *testing.T) {
					testInvalidCoords(
						t,
						locationService,
						data.testedMethod,
						data.expectMethodForGetWeatherFromStorage,
						data.expectMethodForGetWeatherFromProvider,
						data.expectMethodForSaveWeather,
						data.expectMethodForSaveUpdatedTimeWeather,
						forecaster,
					)
				})

			t.Run(
				fmt.Sprintf("%s empty coords in request for forecaster: %s", data.weatherType, forecaster),
				func(t *testing.T) {
					testEmptyCoords(
						t,
						locationService,
						data.testedMethod,
						data.expectMethodForGetWeatherFromStorage,
						data.expectMethodForGetWeatherFromProvider,
						data.expectMethodForSaveWeather,
						data.expectMethodForSaveUpdatedTimeWeather,
						forecaster,
					)
				})

			t.Run(
				fmt.Sprintf("%s empty locale in request for forecaster: %s", data.weatherType, forecaster),
				func(t *testing.T) {
					testEmptyLocale(
						t,
						locationService,
						data.testedMethod,
						data.expectMethodForGetWeatherFromStorage,
						data.expectMethodForGetWeatherFromProvider,
						data.expectMethodForSaveWeather,
						data.expectMethodForSaveUpdatedTimeWeather,
						forecaster,
					)
				})
		}

	}
}

type inputData struct {
	weatherType                           string
	returnedWeather                       any
	testedMethod                          func(ctx context.Context, request *dto.WeatherRequest) (any, error)
	expectMethodForGetWeatherFromStorage  func(ctx, addressHash, forecaster any) *gomock.Call
	expectMethodForGetWeatherFromProvider func(ctx, request any) *gomock.Call
	expectMethodForSaveWeather            func(ctx, addressHash, forecaster, weather any) *gomock.Call
	expectMethodForSaveUpdatedTimeWeather func(ctx, addressHash, forecaster, lastTime any) *gomock.Call
}

func getInputDatas(
	weatherService service.WeatherService,
	weatherStorage *mock_storage.MockWeatherStorage,
	weatherProvider *mock_provider.MockWeatherProvider) *[]inputData {
	return &[]inputData{
		{
			weatherType:                           "Current",
			returnedWeather:                       &dto.CurrentWeather{EpochTime: time.Now().UnixMilli()},
			expectMethodForGetWeatherFromStorage:  weatherStorage.EXPECT().GetCurrentWeather,
			expectMethodForGetWeatherFromProvider: weatherProvider.EXPECT().CurrentWeather,
			expectMethodForSaveWeather:            weatherStorage.EXPECT().SaveCurrentWeather,
			expectMethodForSaveUpdatedTimeWeather: weatherStorage.EXPECT().SaveUpdatedTimeCurrentWeather,
			testedMethod: func(ctx context.Context, request *dto.WeatherRequest) (any, error) {
				return weatherService.CurrentWeather(ctx, request)
			},
		},
		{
			weatherType:                           "Hourly",
			returnedWeather:                       &[]*dto.HourlyWeather{{EpochTime: time.Now().UnixMilli()}},
			expectMethodForGetWeatherFromStorage:  weatherStorage.EXPECT().GetHourlyWeather,
			expectMethodForGetWeatherFromProvider: weatherProvider.EXPECT().HourlyWeather,
			expectMethodForSaveWeather:            weatherStorage.EXPECT().SaveHourlyWeather,
			expectMethodForSaveUpdatedTimeWeather: weatherStorage.EXPECT().SaveUpdatedTimeHourlyWeather,
			testedMethod: func(ctx context.Context, request *dto.WeatherRequest) (any, error) {
				return weatherService.HourlyWeather(ctx, request)
			},
		},
		{
			weatherType:                           "Daily",
			returnedWeather:                       &[]*dto.DailyWeather{{EpochTime: time.Now().UnixMilli()}},
			expectMethodForGetWeatherFromStorage:  weatherStorage.EXPECT().GetDailyWeather,
			expectMethodForGetWeatherFromProvider: weatherProvider.EXPECT().DailyWeather,
			expectMethodForSaveWeather:            weatherStorage.EXPECT().SaveDailyWeather,
			expectMethodForSaveUpdatedTimeWeather: weatherStorage.EXPECT().SaveUpdatedTimeDailyWeather,
			testedMethod: func(ctx context.Context, request *dto.WeatherRequest) (any, error) {
				return weatherService.DailyWeather(ctx, request)
			},
		},
	}
}

func testWeatherNotFoundInStorage(
	t *testing.T,
	locationService *mock_service.MockLocationService,
	testedMethodGetWeatherFromService func(ctx context.Context, request *dto.WeatherRequest) (any, error),
	expectMethodForGetWeatherFromStorage func(ctx, addressHash, forecaster any) *gomock.Call,
	expectMethodForGetWeatherFromProvider func(ctx, request any) *gomock.Call,
	expectMethodForSaveWeather func(ctx, addressHash, forecaster, weather any) *gomock.Call,
	expectMethodForSaveUpdatedTimeWeather func(ctx, addressHash, forecaster, lastTime any) *gomock.Call,
	forecaster dto.WeatherForecaster,
	returnedWeatherFromProvider any,
) {
	ctx := context.Background()
	req := getRequest()
	locationInfo := getLocationInfo()
	locationService.EXPECT().LocationByCoords(gomock.Any(), req.Coords).Return(locationInfo, nil).Times(1)
	expectMethodForGetWeatherFromStorage(gomock.Any(), locationInfo.AddressHash, forecaster).Return(nil, storage.ErrNotFound).Times(1)
	expectMethodForGetWeatherFromProvider(gomock.Any(), gomock.Any()).Return(returnedWeatherFromProvider, nil).Times(1)
	expectMethodForSaveWeather(gomock.Any(), locationInfo.AddressHash, forecaster, returnedWeatherFromProvider).Times(1)
	expectMethodForSaveUpdatedTimeWeather(gomock.Any(), locationInfo.AddressHash, forecaster, gomock.Not(0)).Times(1)

	resp, err := testedMethodGetWeatherFromService(ctx, req)

	assert.Nil(t, err)
	assert.Equal(t, returnedWeatherFromProvider, resp)
}

func testWeatherInStorage(
	t *testing.T,
	locationService *mock_service.MockLocationService,
	testedMethodGetWeatherFromService func(ctx context.Context, request *dto.WeatherRequest) (any, error),
	expectMethodForGetWeatherFromStorage func(ctx, addressHash, forecaster any) *gomock.Call,
	expectMethodForGetWeatherFromProvider func(ctx, request any) *gomock.Call,
	expectMethodForSaveWeather func(ctx, addressHash, forecaster, weather any) *gomock.Call,
	expectMethodForSaveUpdatedTimeWeather func(ctx, addressHash, forecaster, lastTime any) *gomock.Call,
	forecaster dto.WeatherForecaster,
	returnedWeatherDto any) {

	ctx := context.Background()
	req := getRequest()
	locationInfo := getLocationInfo()
	locationService.EXPECT().LocationByCoords(gomock.Any(), req.Coords).Return(locationInfo, nil).Times(1)
	expectMethodForGetWeatherFromStorage(gomock.Any(), locationInfo.AddressHash, forecaster).Return(returnedWeatherDto, nil).Times(1)
	expectMethodForGetWeatherFromProvider(gomock.Any(), gomock.Any()).Times(0)
	expectMethodForSaveWeather(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	expectMethodForSaveUpdatedTimeWeather(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Not(0)).Times(0)

	resp, err := testedMethodGetWeatherFromService(ctx, req)

	assert.Nil(t, err)
	assert.EqualValues(t, returnedWeatherDto, resp)
}

func testWeatherProviderForecasterReturnError(
	t *testing.T,
	locationService *mock_service.MockLocationService,
	testedMethodGetWeatherFromService func(ctx context.Context, request *dto.WeatherRequest) (any, error),
	expectMethodForGetWeatherFromStorage func(ctx, addressHash, forecaster any) *gomock.Call,
	expectMethodForGetWeatherFromProvider func(ctx, request any) *gomock.Call,
	expectMethodForSaveWeather func(ctx, addressHash, forecaster, weather any) *gomock.Call,
	expectMethodForSaveUpdatedTimeWeather func(ctx, addressHash, forecaster, lastTime any) *gomock.Call,
	forecaster dto.WeatherForecaster) {

	expectedErr := http.ErrCountRequestIsOut
	ctx := context.Background()
	req := getRequest()
	locationInfo := getLocationInfo()
	locationService.EXPECT().LocationByCoords(gomock.Any(), req.Coords).Return(locationInfo, nil).Times(1)
	expectMethodForGetWeatherFromStorage(gomock.Any(), locationInfo.AddressHash, forecaster).Return(nil, storage.ErrNotFound).Times(1)
	expectMethodForGetWeatherFromProvider(gomock.Any(), gomock.Any()).Return(nil, expectedErr).Times(1)
	expectMethodForSaveWeather(gomock.Any(), locationInfo.AddressHash, forecaster, gomock.Any()).Times(0)
	expectMethodForSaveUpdatedTimeWeather(gomock.Any(), locationInfo.AddressHash, forecaster, gomock.Nil()).Times(0)

	resp, err := testedMethodGetWeatherFromService(ctx, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, expectedErr)
}

func testInvalidCoords(
	t *testing.T,
	locationService *mock_service.MockLocationService,
	testedMethodGetWeatherFromService func(ctx context.Context, request *dto.WeatherRequest) (any, error),
	expectMethodForGetWeatherFromStorage func(ctx, addressHash, forecaster any) *gomock.Call,
	expectMethodForGetWeatherFromProvider func(ctx, request any) *gomock.Call,
	expectMethodForSaveWeather func(ctx, addressHash, forecaster, weather any) *gomock.Call,
	expectMethodForSaveUpdatedTimeWeather func(ctx, addressHash, forecaster, lastTime any) *gomock.Call,
	forecaster dto.WeatherForecaster) {

	expectedErr := service.ErrInvalidCoords
	req := &dto.WeatherRequest{
		Coords: &dto.Coords{
			Latitude:  -99,
			Longitude: 180,
		},
		Locale: "us",
		Unit:   dto.UnitImperial,
	}

	testInvalidRequest(
		t,
		req,
		expectedErr,
		locationService,
		testedMethodGetWeatherFromService,
		expectMethodForGetWeatherFromStorage,
		expectMethodForGetWeatherFromProvider,
		expectMethodForSaveWeather,
		expectMethodForSaveUpdatedTimeWeather,
		forecaster)
}

func testEmptyCoords(
	t *testing.T,
	locationService *mock_service.MockLocationService,
	testedMethodGetWeatherFromService func(ctx context.Context, request *dto.WeatherRequest) (any, error),
	expectMethodForGetWeatherFromStorage func(ctx, addressHash, forecaster any) *gomock.Call,
	expectMethodForGetWeatherFromProvider func(ctx, request any) *gomock.Call,
	expectMethodForSaveWeather func(ctx, addressHash, forecaster, weather any) *gomock.Call,
	expectMethodForSaveUpdatedTimeWeather func(ctx, addressHash, forecaster, lastTime any) *gomock.Call,
	forecaster dto.WeatherForecaster) {

	expectedErr := service.ErrEmptyCoords
	req := getRequest()
	req.Coords = nil

	testInvalidRequest(
		t,
		req,
		expectedErr,
		locationService,
		testedMethodGetWeatherFromService,
		expectMethodForGetWeatherFromStorage,
		expectMethodForGetWeatherFromProvider,
		expectMethodForSaveWeather,
		expectMethodForSaveUpdatedTimeWeather,
		forecaster)

}

func testEmptyLocale(
	t *testing.T,
	locationService *mock_service.MockLocationService,
	testedMethodGetWeatherFromService func(ctx context.Context, request *dto.WeatherRequest) (any, error),
	expectMethodForGetWeatherFromStorage func(ctx, addressHash, forecaster any) *gomock.Call,
	expectMethodForGetWeatherFromProvider func(ctx, request any) *gomock.Call,
	expectMethodForSaveWeather func(ctx, addressHash, forecaster, weather any) *gomock.Call,
	expectMethodForSaveUpdatedTimeWeather func(ctx, addressHash, forecaster, lastTime any) *gomock.Call,
	forecaster dto.WeatherForecaster) {

	expectedErr := service.ErrEmptyLocale
	req := getRequest()
	req.Locale = ""

	testInvalidRequest(
		t,
		req,
		expectedErr,
		locationService,
		testedMethodGetWeatherFromService,
		expectMethodForGetWeatherFromStorage,
		expectMethodForGetWeatherFromProvider,
		expectMethodForSaveWeather,
		expectMethodForSaveUpdatedTimeWeather,
		forecaster)
}

func testInvalidRequest(
	t *testing.T,
	req *dto.WeatherRequest,
	expectedErr error,
	locationService *mock_service.MockLocationService,
	testedMethodGetWeatherFromService func(ctx context.Context, request *dto.WeatherRequest) (any, error),
	expectMethodForGetWeatherFromStorage func(ctx, addressHash, forecaster any) *gomock.Call,
	expectMethodForGetWeatherFromProvider func(ctx, request any) *gomock.Call, expectMethodForSaveWeather func(ctx, addressHash, forecaster, weather any) *gomock.Call,
	expectMethodForSaveUpdatedTimeWeather func(ctx, addressHash, forecaster, lastTime any) *gomock.Call,
	forecaster dto.WeatherForecaster) {

	ctx := context.Background()

	locationService.EXPECT().LocationByCoords(gomock.Any(), req.Coords).Times(0)
	expectMethodForGetWeatherFromStorage(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	expectMethodForGetWeatherFromProvider(gomock.Any(), gomock.Any()).Times(0)
	expectMethodForSaveWeather(gomock.Any(), gomock.Any(), forecaster, gomock.Any()).Times(0)
	expectMethodForSaveUpdatedTimeWeather(gomock.Any(), gomock.Any(), forecaster, gomock.Nil()).Times(0)

	resp, err := testedMethodGetWeatherFromService(ctx, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, expectedErr)
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
		Coords: getCoords(),
		Locale: "us",
		Unit:   dto.UnitImperial,
	}
}

func getLocationInfo() *dto.LocationInfo {
	return &dto.LocationInfo{
		Coords:                 *getCoords(),
		AddressHash:            "123321",
		AccuWeatherLocationKey: "123",
	}
}

func getCoords() *dto.Coords {
	return &dto.Coords{
		Latitude:  9.11,
		Longitude: 9.11,
	}
}
