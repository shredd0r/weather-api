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
	location, err := s.locationService.LocationByCoords(ctx, request.Coords)
	if err != nil {
		return nil, err
	}

	r := dto.WeatherRequestProviderDto{
		Location: *location,
		Language: request.Language,
		Unit:     request.Unit,
	}

	weather, err := s.weatherStorage.GetCurrentWeatherBy(ctx, r.Location.AddressHash, s.forecaster)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			weather, err = s.weatherProvider.CurrentWeather(ctx, r)
			if err != nil {
				return nil, err
			}
			go func() {
				err := s.weatherStorage.SaveCurrentWeather(ctx, r.Location.AddressHash, s.forecaster, weather)
				if err != nil {
					s.logger.Errorf("error when try save current weather, error: %s", err.Error())
				}
			}()
		} else {
			return nil, err
		}
	}
	return weather, nil
}
func (s *WeatherServiceImpl) HourlyWeather(ctx context.Context, request dto.WeatherRequestDto) (*[]*dto.HourlyWeather, error) {
	return nil, nil
}

func (s *WeatherServiceImpl) DailyWeather(ctx context.Context, request dto.WeatherRequestDto) (*[]*dto.DailyWeather, error) {
	return nil, nil
}
