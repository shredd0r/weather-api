package service

import (
	"context"
	"sync"
	"weather-api/client/http"
	"weather-api/dto"
	"weather-api/log"
	"weather-api/provider"
	"weather-api/storage"
)

type LocationService interface {
	LocationByCoords(ctx context.Context, coords *dto.Coords) (*dto.LocationInfo, error)
}

type LocationServiceImpl struct {
	wg               *sync.WaitGroup
	locationProvider provider.LocationProvider
}

func NewLocationService(logger log.Logger, wg *sync.WaitGroup, locationStorage storage.LocationStorage,
	accuWeatherClient http.AccuWeatherInterface, apiNinjasClient http.ApiNinjasInterface) LocationService {
	return &LocationServiceImpl{
		locationProvider: provider.NewLocationProvider(logger, wg, locationStorage, accuWeatherClient, apiNinjasClient),
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
