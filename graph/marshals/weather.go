package marshals

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/shredd0r/weather-api/dto"
)

func MarshalCurrentWeather(v dto.CurrentWeather) graphql.Marshaler {
	return graphql.MarshalAny(v)
}

func UnmarshalCurrentWeather(v any) (dto.CurrentWeather, error) {
	return v.(dto.CurrentWeather), nil
}

func MarshalHourlyWeather(v dto.HourlyWeather) graphql.Marshaler {
	return graphql.MarshalAny(v)
}

func UnmarshalHourlyWeather(v any) (dto.HourlyWeather, error) {
	return v.(dto.HourlyWeather), nil
}

func MarshalDailyWeather(v dto.DailyWeather) graphql.Marshaler {
	return graphql.MarshalAny(v)
}

func UnmarshalDailyWeather(v any) (dto.DailyWeather, error) {
	return v.(dto.DailyWeather), nil
}
