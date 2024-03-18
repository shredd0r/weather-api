package util

import (
	"fmt"
	"math"
	"weather-api/dto"
	"weather-api/entity"
)

func MapToCurrentWeatherDto(entity *entity.CurrentWeatherEntity) dto.CurrentWeatherDto {
	return dto.CurrentWeatherDto{
		EpochTime:            entity.EpochTime,
		Visibility:           entity.Visibility,
		CurrentTemperature:   entity.CurrentTemperature,
		MinTemperature:       &entity.MinTemperature,
		MaxTemperature:       &entity.MaxTemperature,
		FeelsLikeTemperature: entity.FillLikeTemperature,
		IconResource:         &entity.IconResource,
	}
}

func MapToHourlyWeatherDto(entity entity.HourlyWeatherEntity) dto.HourlyWeatherDto {
	return dto.HourlyWeatherDto{}
}

func MapToLocalityDtoBy(entity *entity.LocalityEntity, cityName string) dto.LocalityDto {
	return dto.LocalityDto{
		Latitude:               entity.Latitude,
		Longitude:              entity.Longitude,
		AccuWeatherLocationKey: entity.AccuWeatherLocationKey,
		CityName:               cityName,
	}
}

func Float64ToString(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func RoundFloat64(f float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(f*ratio) / ratio
}

func GetIconResourceNameByAccuWeatherIcon(iconId uint8) *string {
	return nil
}
