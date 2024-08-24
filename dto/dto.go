package dto

import (
	"encoding/json"
	"strings"

	"github.com/shredd0r/weather-api/util"
)

type IconId string

const (
	IconIdUndefined   = "0"
	IconIdClearDay    = "1"
	IconIdClearNight  = "2"
	IconIdCloudyDay   = "3"
	IconIdCloudyNight = "4"
	IconIdCloudy      = "5"
	IconIdRainyDay    = "6"
	IconIdRainyNight  = "7"
	IconIdRainy       = "8"
	IconIdSnowyDay    = "9"
	IconIdSnowyNight  = "10"
	IconIdSnowy       = "11"
	IconIdThunder     = "12"
	IconIdMist        = "13"
)

type Error struct {
	Message string
}

type CurrentWeather struct {
	EpochTime            int64    `json:"epochTime"`
	Visibility           *float64 `json:"visibility"`
	CurrentTemperature   *float64 `json:"currentTemperature"`
	MinTemperature       *float64 `json:"minTemperature"`
	MaxTemperature       *float64 `json:"maxTemperature"`
	FeelsLikeTemperature *float64 `json:"feelsLikeTemperature"`
	IconId               *IconId  `json:"iconId"`
	MobileLink           string   `json:"mobileLink"`
	Link                 string   `json:"link"`
}

func (w *CurrentWeather) String() string {
	return getJSONStr(w)
}

type HourlyWeather struct {
	EpochTime                  int64             `json:"epochTime"`
	CurrentTemperature         *float64          `json:"currentTemperature"`
	FeelsLikeTemperature       *float64          `json:"feelsLikeTemperature"`
	UVIndex                    *float64          `json:"uvIndex"`
	ProbabilityOfPrecipitation *float64          `json:"probabilityOfPrecipitation"`
	PrecipitationType          PrecipitationType `json:"precipitationType"`
	AmountOfPrecipitation      *float64          `json:"amountOfPrecipitation"`
	Wind                       *Wind             `json:"wind"`
	IconId                     *IconId           `json:"iconId"`
	MobileLink                 string            `json:"mobileLink"`
	Link                       string            `json:"link"`
}

func (w *HourlyWeather) String() string {
	return getJSONStr(w)
}

type DailyWeather struct {
	EpochTime                  int64             `json:"epochTime"`
	MinTemperature             *float64          `json:"minTemperature"`
	MaxTemperature             *float64          `json:"maxTemperature"`
	Humidity                   float64           `json:"humidity"`
	UVIndex                    *float64          `json:"uvIndex"`
	SunriseTime                int64             `json:"sunriseTime"`
	SunsetTime                 int64             `json:"sunsetTime"`
	Wind                       *Wind             `json:"wind"`
	ProbabilityOfPrecipitation *float64          `json:"probabilityOfPrecipitation"`
	PrecipitationType          PrecipitationType `json:"precipitationType"`
	IconId                     *IconId           `json:"iconId"`
	MobileLink                 string            `json:"mobileLink"`
	Link                       string            `json:"link"`
}

func (w *DailyWeather) String() string {
	return getJSONStr(w)
}

type Wind struct {
	Speed   *float64 `json:"speed"`
	Degrees float64  `json:"degrees"`
}

type Coords struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (w *Coords) String() string {
	return getJSONStr(w)
}

func NewCoords(strCoords string) Coords {
	coords := strings.Split(strCoords, ",")
	if len(coords) != 2 {
		return Coords{}
	}

	return Coords{
		Latitude:  util.StringToFloat64(coords[0]),
		Longitude: util.StringToFloat64(coords[1]),
	}
}

type Location struct {
	Coords                 Coords `json:"coords"`
	AccuWeatherLocationKey string `json:"accuWeatherLocationKey"`
}

type LocationInfo struct {
	Coords                 Coords
	AddressHash            string
	AccuWeatherLocationKey string
}

type GeocodingRequest struct {
	City    *string
	State   *string
	Country *string
}

type Geocoding struct {
	Name      string
	Latitude  float64
	Longitude float64
	Country   string
	State     string
}

type ReverseGeocodingRequest struct {
	Latitude  float64
	Longitude float64
}

type ReverseGeocoding struct {
	Country string
	Name    string
	State   string
}

type GeoPositionResponse struct {
	Version int
	Key     string
}

type WeatherRequest struct {
	Coords *Coords
	Locale string
	Unit
}

type WeatherRequestProvider struct {
	Location LocationInfo
	Locale   string
	Unit
}

type WeatherForecaster string

const (
	WeatherForecasterUnspecified = ""
	WeatherForecasterAccuWeather = "AccuWeather"
	WeatherForecasterOpenWeather = "OpenWeather"
	WeatherForecasterWeatherApi  = "WeatherApi"
)

type Unit string

const (
	UnitUnspecified = ""
	UnitImperial    = "imperial"
	UnitMetric      = "metric"
)

type PrecipitationType string

const (
	PrecipitationTypeEmpty = ""
	PrecipitationTypeRain  = "Rain"
	PrecipitationTypeSnow  = "Snow"
	PrecipitationTypeIce   = "Ice"
)

func getJSONStr(obj any) string {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}
