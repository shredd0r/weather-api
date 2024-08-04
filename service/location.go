package service

import (
	"context"
	"github.com/shredd0r/weather-api/client/http"
	"github.com/shredd0r/weather-api/dto"
	"github.com/shredd0r/weather-api/log"
	"github.com/shredd0r/weather-api/provider"
	"github.com/shredd0r/weather-api/storage"
	"sync"
)

type LocationService interface {
	LocationByCoords(ctx context.Context, coords *dto.Coords) (*dto.LocationInfo, error)
}

type LocationServiceImpl struct {
	wg               *sync.WaitGroup
	locationProvider provider.LocationProvider
}

func NewLocationService(logger log.Logger, locationStorage storage.LocationStorage,
	accuWeatherClient http.AccuWeatherInterface, apiNinjasClient http.ApiNinjasInterface) LocationService {
	return &LocationServiceImpl{
		locationProvider: provider.NewLocationProvider(logger, locationStorage, accuWeatherClient, apiNinjasClient),
	}
}

func (s *LocationServiceImpl) LocationByCoords(ctx context.Context, coords *dto.Coords) (*dto.LocationInfo, error) {
	addressHash, err := s.locationProvider.GetAddressHashByCoords(ctx, coords)
	if err != nil {
		return nil, err
	}
	location, err := s.locationProvider.LocationByAddressHash(ctx, coords, addressHash)
	if err != nil {
		return nil, err
	}
	return &dto.LocationInfo{
		Coords:                 location.Coords,
		AddressHash:            addressHash,
		AccuWeatherLocationKey: location.AccuWeatherLocationKey,
	}, nil
}
