package provider

import "weather-api/dto"

type LocalityProvider interface {
	GetLocalityDto(cityName string) (*dto.LocalityDto, error)
	GetCityName(latitude float64, longitude float64) (*string, error)
}

type WeatherProvider interface {
	CurrentWeatherInfo(weatherRequestDto dto.WeatherRequestDto) *dto.CurrentWeatherDto
	HourlyWeatherInfo(weatherRequestDto dto.WeatherRequestDto) *[]dto.HourlyWeatherDto
	DailyWeatherInfo(weatherRequestDto dto.WeatherRequestDto) *[]dto.DailyWeatherDto
}

type SaveInfoProvider[T any] interface {
	Save(info T) error
}
