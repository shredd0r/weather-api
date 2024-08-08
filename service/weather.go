package service

import (
	"context"
	"errors"
	"github.com/shredd0r/weather-api/dto"
	"github.com/shredd0r/weather-api/log"
	"github.com/shredd0r/weather-api/provider"
	"github.com/shredd0r/weather-api/storage"
	"sync"
	"time"
)

var (
	ErrEmptyCoords   = errors.New("coords is empty")
	ErrEmptyLocale   = errors.New("locale is empty")
	ErrInvalidCoords = errors.New("invalid coords latitude range is [-90, 90], longitude range is [-180, 180)")
)

type WeatherService interface {
	CurrentWeather(ctx context.Context, request *dto.WeatherRequest) (*dto.CurrentWeather, error)
	HourlyWeather(ctx context.Context, request *dto.WeatherRequest) (*[]*dto.HourlyWeather, error)
	DailyWeather(ctx context.Context, request *dto.WeatherRequest) (*[]*dto.DailyWeather, error)
}

type WeatherServiceImpl struct {
	wg              *sync.WaitGroup
	forecaster      dto.WeatherForecaster
	logger          log.Logger
	locationService LocationService
	weatherProvider provider.WeatherProvider
	weatherStorage  storage.WeatherStorage
}

func NewWeatherService(logger log.Logger, forecaster dto.WeatherForecaster, locationService LocationService,
	weatherProvider provider.WeatherProvider, weatherStorage storage.WeatherStorage) WeatherService {
	return &WeatherServiceImpl{
		logger:          logger,
		wg:              &sync.WaitGroup{},
		forecaster:      forecaster,
		locationService: locationService,
		weatherProvider: weatherProvider,
		weatherStorage:  weatherStorage,
	}
}

func (s *WeatherServiceImpl) CurrentWeather(ctx context.Context, request *dto.WeatherRequest) (*dto.CurrentWeather, error) {
	return workflowGetWeatherFromProviderOrStorage[dto.CurrentWeather](
		s.logger,
		ctx,
		request,
		s.forecaster,
		s.getWeatherRequestProviderDto,
		s.weatherStorage.GetCurrentWeather,
		s.weatherProvider.CurrentWeather,
		s.weatherStorage.SaveUpdatedTimeCurrentWeather,
		s.weatherStorage.SaveCurrentWeather)
}
func (s *WeatherServiceImpl) HourlyWeather(ctx context.Context, request *dto.WeatherRequest) (*[]*dto.HourlyWeather, error) {
	return workflowGetWeatherFromProviderOrStorage[[]*dto.HourlyWeather](
		s.logger,
		ctx,
		request,
		s.forecaster,
		s.getWeatherRequestProviderDto,
		s.weatherStorage.GetHourlyWeather,
		s.weatherProvider.HourlyWeather,
		s.weatherStorage.SaveUpdatedTimeHourlyWeather,
		s.weatherStorage.SaveHourlyWeather)
}

func (s *WeatherServiceImpl) DailyWeather(ctx context.Context, request *dto.WeatherRequest) (*[]*dto.DailyWeather, error) {
	return workflowGetWeatherFromProviderOrStorage[[]*dto.DailyWeather](
		s.logger,
		ctx,
		request,
		s.forecaster,
		s.getWeatherRequestProviderDto,
		s.weatherStorage.GetDailyWeather,
		s.weatherProvider.DailyWeather,
		s.weatherStorage.SaveUpdatedTimeDailyWeather,
		s.weatherStorage.SaveDailyWeather)
}

func (s *WeatherServiceImpl) getWeatherRequestProviderDto(ctx context.Context, request *dto.WeatherRequest) (*dto.WeatherRequestProvider, error) {
	location, err := s.locationService.LocationByCoords(ctx, request.Coords)
	if err != nil {
		return nil, err
	}

	return &dto.WeatherRequestProvider{
		Location: *location,
		Locale:   request.Locale,
		Unit:     request.Unit,
	}, nil
}

func workflowGetWeatherFromProviderOrStorage[T any](
	logger log.Logger, ctx context.Context, request *dto.WeatherRequest, forecaster dto.WeatherForecaster,
	methodForGetWeatherProviderRequest func(ctx context.Context, request *dto.WeatherRequest) (*dto.WeatherRequestProvider, error),
	methodForGetFromStorage func(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (*T, error),
	methodForGetFromProvider func(ctx context.Context, request *dto.WeatherRequestProvider) (*T, error),
	methodForSaveUpdatedTime func(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, lastTime int64) error,
	methodForSaveWeather func(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather T) error) (*T, error) {

	ok, err := validateRequest(request)
	if !ok {
		return nil, err
	}

	wg := &sync.WaitGroup{}
	r, err := methodForGetWeatherProviderRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	weather, err := methodForGetFromStorage(ctx, r.Location.AddressHash, forecaster)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			updatedTime := time.Now().UnixMilli()
			weather, err = methodForGetFromProvider(ctx, r)
			if err != nil {
				return nil, err
			}
			wg.Add(2)
			go func() {
				defer wg.Done()

				logger.Info("start save new weather info")
				err := methodForSaveWeather(ctx, r.Location.AddressHash, forecaster, *weather)
				if err != nil {
					logger.Errorf("error when try save daily weather, error: %s", err.Error())
				}
			}()
			go func() {
				defer wg.Done()

				logger.Info("start save updated time for weather")
				err := methodForSaveUpdatedTime(ctx, r.Location.AddressHash, forecaster, updatedTime)
				if err != nil {
					logger.Errorf("error when try save daily weather, error: %s", err.Error())
				}
			}()

		} else {
			return nil, err
		}
	}

	wg.Wait()
	return weather, nil
}

func validateRequest(request *dto.WeatherRequest) (bool, error) {
	ok, err := coordsIsCorrect(request)
	if !ok {
		return false, err
	}

	if request.Locale == "" {
		return false, ErrEmptyLocale
	}

	return true, nil
}

func coordsIsCorrect(request *dto.WeatherRequest) (bool, error) {
	if request.Coords != nil {
		if latitudeIsCorrect(request.Coords) && longitudeIsCorrect(request.Coords) {
			return true, nil
		}
		return false, ErrInvalidCoords
	}
	return false, ErrEmptyCoords
}

func latitudeIsCorrect(coords *dto.Coords) bool {
	return coords.Latitude <= 90 || coords.Latitude >= -90
}

func longitudeIsCorrect(coords *dto.Coords) bool {
	return coords.Longitude < 180 || coords.Longitude >= -180
}
