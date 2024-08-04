// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"github.com/shredd0r/weather-api/dto"
)

type Coords struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Query struct {
}

type WeatherRequest struct {
	Coords     *Coords               `json:"coords"`
	Locale     string                `json:"locale"`
	Unit       dto.Unit              `json:"unit"`
	Forecaster dto.WeatherForecaster `json:"forecaster"`
}
