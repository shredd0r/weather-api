package provider

import (
	"github.com/sirupsen/logrus"
	"time"
	"weather-api/client/http"
	"weather-api/dto"
	"weather-api/entity"
	"weather-api/repository"
	"weather-api/util"
)

type HttpLocalityProvider struct {
	log               *logrus.Logger
	accuWeatherClient *http.AccuWeatherClient
	ninjaClient       *http.ApiNinjasClient
}

func NewHttpLocalityProvider(log *logrus.Logger, accuWeatherClient *http.AccuWeatherClient, ninjaClient *http.ApiNinjasClient) *HttpLocalityProvider {
	return &HttpLocalityProvider{
		log:               log,
		accuWeatherClient: accuWeatherClient,
		ninjaClient:       ninjaClient,
	}
}

func (p *HttpLocalityProvider) GetLocalityDto(cityName string) (*dto.LocalityDto, error) {
	ninjaResp, err := p.ninjaGeocodingRequest(cityName)

	if err != nil {
		return nil, err
	}

	accuWeatherResp, err := p.accuWeatherGeocodingRequest(ninjaResp.Latitude, ninjaResp.Longitude)

	if err != nil {
		return nil, err
	}

	return &dto.LocalityDto{
		Latitude:               ninjaResp.Latitude,
		Longitude:              ninjaResp.Longitude,
		CityName:               cityName,
		AccuWeatherLocationKey: accuWeatherResp.Key,
	}, nil
}

func (p *HttpLocalityProvider) GetCityName(latitude float64, longitude float64) (*string, error) {
	resp, err := p.ninjaClient.GetReversGeocoding(dto.ApiNinjasReverseGeocodingRequestDto{
		Latitude:  latitude,
		Longitude: longitude,
	})

	if err != nil {
		p.log.Fatalf("error from reverse geocoding request %s", err)
		return nil, err
	}

	return &resp.Name, nil
}

type CacheLocalityProvider struct {
	log         *logrus.Logger
	localityRep repository.LocalityRepository
}

func (p *CacheLocalityProvider) GetLocalityDto(cityName string) (*dto.LocalityDto, error) {
	localityEntity, err := p.localityRep.FindByCityName(cityName)

	if err == nil {
		localityDto := util.MapToLocalityDtoBy(localityEntity, cityName)
		return &localityDto, nil
	}

	return nil, err
}

func (p *CacheLocalityProvider) Save(localityDto dto.LocalityDto) error {
	err := p.localityRep.Save(entity.LocalityEntity{
		Latitude:               localityDto.Latitude,
		Longitude:              localityDto.Longitude,
		AccuWeatherLocationKey: localityDto.AccuWeatherLocationKey,
		AddedEpochTime:         time.Now().Unix(),
	})

	return err
}

func (p *CacheLocalityProvider) GetCityName(latitude float64, longitude float64) (*string, error) {
	return nil, repository.ErrNotFound
}

func (p *HttpLocalityProvider) ninjaGeocodingRequest(cityName string) (*dto.ApiNinjasGeocodingResponseDto, error) {
	resp, err := p.ninjaClient.GetGeocoding(
		dto.ApiNinjasGeocodingRequestDto{
			City: cityName,
		})

	if err != nil {
		p.log.Fatalf("error from geocoding request %s", err)
		return nil, err
	}

	return resp, nil
}

func (p *HttpLocalityProvider) accuWeatherGeocodingRequest(latitude float64, longitude float64) (*dto.AccuWeatherGeoPositionResponseDto, error) {
	resp, err := p.accuWeatherClient.GetGeoPositionSearch(
		dto.AccuWeatherGeoPositionRequestDto{
			Latitude:  latitude,
			Longitude: longitude,
		})

	if err != nil {
		p.log.Fatalf("error from geo position request %s", err)
		return nil, err
	}

	return resp, nil
}
