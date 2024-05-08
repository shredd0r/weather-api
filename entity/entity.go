package entity

import (
	"github.com/google/uuid"
	"weather-api/dto"
)

type Entity interface {
	getId() uuid.UUID
}

type CurrentWeatherEntity struct {
	Id                   uuid.UUID
	Locality             *LocalityEntity
	Forecaster           dto.WeatherForecaster
	EpochTime            int64
	Visibility           *float64
	CurrentTemperature   *float64
	MinTemperature       *float64
	MaxTemperature       *float64
	FeelsLikeTemperature *float64
	IconResource         *string
	MobileLink           string
	Link                 string
	AddedEpochTime       int64
}

func (e CurrentWeatherEntity) getId() uuid.UUID {
	return e.Id
}

type HourlyWeatherEntity struct {
	Id                         uuid.UUID
	Locality                   *LocalityEntity
	Forecaster                 dto.WeatherForecaster
	Temperature                *float64
	FeelsLikeTemperature       *float64
	UVIndex                    *uint8
	EpochTime                  int64
	ProbabilityOfPrecipitation *float64
	PrecipitationType          dto.PrecipitationType
	AmountOfPrecipitation      *float64
	Wind                       WindEntity
	IconResource               *string
	MobileLink                 string
	Link                       string
	AddedEpochTime             int64
}

func (e HourlyWeatherEntity) getId() uuid.UUID {
	return e.Id
}

type DailyWeatherEntity struct {
	Id                         uuid.UUID
	Locality                   *LocalityEntity
	Forecaster                 dto.WeatherForecaster
	EpochTime                  int64
	MinTemperature             *float64
	MaxTemperature             *float64
	Humidity                   float64
	UVIndex                    *float64
	SunriseTime                int64
	SunsetTime                 int64
	Wind                       *WindEntity
	ProbabilityOfPrecipitation *float64
	PrecipitationType          dto.PrecipitationType
	IconResource               *string
	MobileLink                 string
	Link                       string
	AddedEpochTime             int64
}

func (e DailyWeatherEntity) getId() uuid.UUID {
	return e.Id
}

type WindEntity struct {
	Speed   *float64
	Degrees float64
}

type LocalityEntity struct {
	Id                     uuid.UUID
	Latitude               float64
	Longitude              float64
	AccuWeatherLocationKey string
	CityName               string
	AddedEpochTime         int64
}

func (e LocalityEntity) getId() uuid.UUID {
	return e.Id
}
