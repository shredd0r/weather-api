package provider

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"hash"
	"sync"
	"time"
	"weather-api/client/http"
	"weather-api/dto"
	"weather-api/log"
	"weather-api/storage"
)

type LocationProviderImpl struct {
	wg                *sync.WaitGroup
	hasher            hash.Hash
	logger            log.Logger
	locationStorage   storage.LocationStorage
	accuWeatherClient http.AccuWeatherInterface
	apiNinjaClient    http.ApiNinjasInterface
}

func NewLocationProvider(logger log.Logger, wg *sync.WaitGroup, locationStorage storage.LocationStorage,
	accuWeatherClient http.AccuWeatherInterface, apiNinjasClient http.ApiNinjasInterface) LocationProvider {
	return &LocationProviderImpl{
		wg:                wg,
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
	p.logger.Info("start get location by address hash")

	location, err := p.locationStorage.GetLocation(ctx, addressHash)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			p.logger.Info("in storage not found location by address hash")
			locationKey, err := p.getLocationKeyFromAccuWeather(coords)
			if err != nil {
				p.logger.Error("error when try get location key from accu weather")
				return nil, err
			}
			location = &dto.Location{
				Coords:                 *coords,
				AccuWeatherLocationKey: locationKey,
			}
			p.storeNewLocation(ctx, location, addressHash)
		}
	}
	return location, nil
}

func (p *LocationProviderImpl) getAddressHashByCoordsFromNinjaApi(coords *dto.Coords) (string, error) {
	p.logger.Info("start get country by coords from ninja api")
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
	lastTime := time.Now().UnixMilli()
	p.logger.Info("start store address hash to location storage")

	p.wg.Add(2)
	go func() {
		defer p.wg.Done()

		err := p.locationStorage.AddCoords(ctx, coords, addressHash)
		if err != nil {
			p.logger.Errorf("error when try new coords in storage, error: %s", err.Error())
		}
	}()
	go func() {
		defer p.wg.Done()

		err := p.locationStorage.UpdateLastTimeUseCoords(ctx, coords, lastTime)
		if err != nil {
			p.logger.Errorf("error when try update last time get address hash, error: %s", err.Error())
		}
	}()

	p.wg.Wait()

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

func (p *LocationProviderImpl) storeNewLocation(ctx context.Context, location *dto.Location, addressHash string) {
	p.logger.Info("start store new location")

	go func() {
		err := p.locationStorage.SaveLocation(ctx, *location, addressHash)
		if err != nil {
			p.logger.Errorf("error when try save new location, error: %s", err.Error())
		}
	}()
}
