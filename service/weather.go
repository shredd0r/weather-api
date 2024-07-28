package service

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"weather-api/dto"
	"weather-api/provider"
	"weather-api/storage"
)

type WeatherService interface {
	CurrentWeather(ctx context.Context, request dto.WeatherRequestDto) (*dto.CurrentWeather, error)
	HourlyWeather(ctx context.Context, request dto.WeatherRequestDto) (*[]*dto.HourlyWeather, error)
	DailyWeather(ctx context.Context, request dto.WeatherRequestDto) (*[]*dto.DailyWeather, error)
}

type WeatherServiceImpl struct {
	forecaster      dto.WeatherForecaster
	logger          *logrus.Logger
	locationService LocationService
	weatherProvider provider.WeatherProvider
	weatherStorage  storage.WeatherStorage
}

func NewWeatherService(logger *logrus.Logger, forecaster dto.WeatherForecaster, locationService LocationService,
	weatherProvider provider.WeatherProvider, weatherStorage storage.WeatherStorage) WeatherService {
	return &WeatherServiceImpl{
		logger:          logger,
		forecaster:      forecaster,
		locationService: locationService,
		weatherProvider: weatherProvider,
		weatherStorage:  weatherStorage,
	}
}

func (s *WeatherServiceImpl) CurrentWeather(ctx context.Context, request dto.WeatherRequestDto) (*dto.CurrentWeather, error) {
	r, err := s.getWeatherRequestProviderDto(ctx, request)
	if err != nil {
		return nil, err
	}

	return workflowGetWeatherFromProviderOrStorage[dto.CurrentWeather](
		s.logger,
		ctx,
		r,
		s.forecaster,
		s.weatherStorage.GetCurrentWeatherBy,
		s.weatherProvider.CurrentWeather,
		s.weatherStorage.SaveCurrentWeather)
}
func (s *WeatherServiceImpl) HourlyWeather(ctx context.Context, request dto.WeatherRequestDto) (*[]*dto.HourlyWeather, error) {
	r, err := s.getWeatherRequestProviderDto(ctx, request)
	if err != nil {
		return nil, err
	}

	return workflowGetWeatherFromProviderOrStorage[[]*dto.HourlyWeather](
		s.logger,
		ctx,
		r,
		s.forecaster,
		s.weatherStorage.GetHourlyWeatherBy,
		s.weatherProvider.HourlyWeather,
		s.weatherStorage.SaveHourlyWeather)
}

func (s *WeatherServiceImpl) DailyWeather(ctx context.Context, request dto.WeatherRequestDto) (*[]*dto.DailyWeather, error) {
	r, err := s.getWeatherRequestProviderDto(ctx, request)
	if err != nil {
		return nil, err
	}

	return workflowGetWeatherFromProviderOrStorage[[]*dto.DailyWeather](
		s.logger,
		ctx,
		r,
		s.forecaster,
		s.weatherStorage.GetDailyWeatherBy,
		s.weatherProvider.DailyWeather,
		s.weatherStorage.SaveDailyWeather)
}

func (s *WeatherServiceImpl) getWeatherRequestProviderDto(ctx context.Context, request dto.WeatherRequestDto) (*dto.WeatherRequestProviderDto, error) {
	location, err := s.locationService.LocationByCoords(ctx, request.Coords)
	if err != nil {
		return nil, err
	}

	return &dto.WeatherRequestProviderDto{
		Location: *location,
		Locale:   request.Locale,
		Unit:     request.Unit,
	}, nil
}

func workflowGetWeatherFromProviderOrStorage[T any](
	logger *logrus.Logger, ctx context.Context, request *dto.WeatherRequestProviderDto, forecaster dto.WeatherForecaster,
	methodForGetFromStorage func(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (*T, error),
	methodForGetFromProvider func(ctx context.Context, request *dto.WeatherRequestProviderDto) (*T, error),
	methodForSaveWeather func(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather T) error) (*T, error) {

	weather, err := methodForGetFromStorage(ctx, request.Location.AddressHash, forecaster)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			weather, err = methodForGetFromProvider(ctx, request)
			if err != nil {
				return nil, err
			}
			go func() {
				err := methodForSaveWeather(ctx, request.Location.AddressHash, forecaster, *weather)
				if err != nil {
					logger.Errorf("error when try save daily weather, error: %s", err.Error())
				}
			}()
		} else {
			return nil, err
		}
	}
	return weather, nil
}
