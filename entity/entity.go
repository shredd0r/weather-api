package entity

import "weather-api/dto"

type CurrentWeatherEntity struct {
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

type HourlyWeatherEntity struct {
	Temperature                *float64
	FeelsLikeTemperature       *float64
	UVIndex                    *int
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

type DailyWeatherEntity struct {
	EpochTime                  int64
	MinTemperature             *float64
	MaxTemperature             *float64
	Humidity                   float64
	UVIndex                    *float64
	SunriseTime                int64
	SunsetTime                 int64
	Wind                       WindEntity
	ProbabilityOfPrecipitation *float64
	PrecipitationType          dto.PrecipitationType
	IconResource               *string
	MobileLink                 string
	Link                       string
	AddedEpochTime             int64
}

type WindEntity struct {
	Speed   *float64
	Degrees float64
}

type LocalityEntity struct {
	Latitude               float64
	Longitude              float64
	AccuWeatherLocationKey string
	AddedEpochTime         int64
}
