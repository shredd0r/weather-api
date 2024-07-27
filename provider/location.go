package provider

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/sirupsen/logrus"
	"hash"
	"weather-api/client/http"
	"weather-api/dto"
	"weather-api/storage"
)

type LocationProviderImpl struct {
	hasher            hash.Hash
	logger            *logrus.Logger
	locationStorage   storage.LocationStorage
	accuWeatherClient http.AccuWeatherInterface
	apiNinjaClient    http.ApiNinjasInterface
}

func NewLocationProvider(logger *logrus.Logger, locationStorage storage.LocationStorage,
	accuWeatherClient http.AccuWeatherInterface, apiNinjasClient http.ApiNinjasInterface) LocationProvider {
	return &LocationProviderImpl{
		hasher:            md5.New(),
		logger:            logger,
		locationStorage:   locationStorage,
		accuWeatherClient: accuWeatherClient,
		apiNinjaClient:    apiNinjasClient,
	}
}

func (p *LocationProviderImpl) GetAddressHashByCoords(ctx context.Context, coords *dto.Coords) (addressHash string, err error) {
	addressHash, err = p.locationStorage.GetAddressHashByCoords(ctx, coords)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			p.logger.Debugf("not found cached coords: %f, %f", coords.Latitude, coords.Longitude)
			addressHash, err = p.getAddressHashByCoordsFromNinjaApi(coords)
			if err != nil {
				return
			}
			p.storeAddressHash(ctx, coords, addressHash)
		}
	}
	return
}

func (p *LocationProviderImpl) LocationByAddressHash(ctx context.Context, coords *dto.Coords, addressHash string) (*dto.Location, error) {
	location, err := p.locationStorage.GetLocation(ctx, addressHash)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			locationKey, err := p.getLocationKeyFromAccuWeather(coords)
			if err != nil {
				return nil, err
			}
			location = &dto.Location{
				Coords:                 *coords,
				AccuWeatherLocationKey: locationKey,
			}
			p.storeNewLocation(ctx, location)
		}
	}
	return location, nil
}

func (p *LocationProviderImpl) getAddressHashByCoordsFromNinjaApi(coords *dto.Coords) (string, error) {
	r, err := p.apiNinjaClient.GetReversGeocoding(dto.ApiNinjasReverseGeocodingRequestDto{
		Latitude:  coords.Latitude,
		Longitude: coords.Longitude,
	})

	if err != nil {
		p.logger.Error("error when try get revers geocoding from api ninja")
		return "", err
	}

	country := (*r)[0]

	addressString := country.Country + country.State + country.Name
	p.hasher.Write([]byte(addressString))
	return hex.EncodeToString(p.hasher.Sum(nil)), nil
}

func (p *LocationProviderImpl) storeAddressHash(ctx context.Context, coords *dto.Coords, addressHash string) {
	go func() {
		err := p.locationStorage.AddNewCoords(ctx, coords, addressHash)
		if err != nil {
			p.logger.Errorf("error when try new coords in storage, error: %s", err.Error())
		}
	}()
	go func() {
		err := p.locationStorage.UpdateLastTimeGetAddressHash(ctx, coords)
		if err != nil {
			p.logger.Errorf("error when try update last time get address hash, error: %s", err.Error())
		}
	}()

	return
}

func (p *LocationProviderImpl) getLocationKeyFromAccuWeather(coords *dto.Coords) (string, error) {
	resp, err := p.accuWeatherClient.GetGeoPositionSearch(dto.AccuWeatherGeoPositionRequestDto{
		Latitude:  coords.Latitude,
		Longitude: coords.Longitude,
	})

	if err != nil {
		return "", err
	}

	return resp.Key, nil
}

func (p *LocationProviderImpl) storeNewLocation(ctx context.Context, location *dto.Location) {
	go func() {
		err := p.locationStorage.SaveLocation(ctx, location)
		if err != nil {
			p.logger.Error("error when try save new location, error: %s", err.Error())
		}
	}()
}
