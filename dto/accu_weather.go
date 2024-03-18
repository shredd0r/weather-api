package dto

type AccuWeatherBaseRequestDto struct {
	Language string
	Metric   bool
	Details  bool
}

type AccuWeatherRequestDto struct {
	AccuWeatherBaseRequestDto
	LocationKey string
}

type AccuWeatherCurrentResponseDto struct {
	EpochTime           int64
	WeatherText         string
	WeatherIcon         uint8
	PrecipitationType   PrecipitationType
	Temperature         AccuWeatherIndicationInfoDto
	RealFeelTemperature AccuWeatherIndicationInfoDto
	UVIndex             uint8
	Visibility          AccuWeatherIndicationInfoDto
	MobileLink          string
	Link                string
}

type AccuWeatherHourlyResponseDto struct {
	EpochTime                int64
	WeatherIcon              uint8
	IconPhrase               string
	Temperature              AccuWeatherValueInfoDto
	RealFeelTemperature      AccuWeatherValueInfoDto
	Wind                     AccuWeatherWindInfoDto
	UVIndex                  uint8
	HasPrecipitation         bool
	PrecipitationType        PrecipitationType
	PrecipitationProbability int
	RainProbability          uint
	SnowProbability          uint
	IceProbability           uint
	TotalLiquid              AccuWeatherValueInfoDto
	Rain                     AccuWeatherValueInfoDto
	Snow                     AccuWeatherValueInfoDto
	Ice                      AccuWeatherValueInfoDto
	MobileLink               string
	Link                     string
}

type AccuWeatherDailyResponseDto struct {
	Headline       AccuWeatherHeadlineDto
	DailyForecasts []AccuWeatherDailyForecastDto
}

type AccuWeatherWindInfoDto struct {
	Speed     AccuWeatherValueInfoDto
	Direction AccuWeatherDirectionDto
}

type AccuWeatherDirectionDto struct {
	Degrees   float64
	Localized string
	English   string
}

type AccuWeatherHeadlineDto struct {
	EffectiveDate      string
	EffectiveEpochDate int64
	Severity           int
	Text               string
	Category           string
	EndDate            string
	EndEpochDate       int64
	MobileLink         string
	Link               string
}

type AccuWeatherDailyForecastDto struct {
	Date        string
	EpochDate   int64
	Temperature AccuWeatherTemperatureDto
	Day         AccuWeatherDayInfoDto
	Night       AccuWeatherDayInfoDto
	Sources     *[]string
	MobileLink  string
	Link        string
}

type AccuWeatherDayInfoDto struct {
	Icon                     uint8
	IconPhrase               string
	HasPrecipitation         bool
	PrecipitationType        PrecipitationType
	PrecipitationProbability int
	Wind                     AccuWeatherWindInfoDto
	TotalLiquid              AccuWeatherValueInfoDto
	RelativeHumidity         AccuWeatherRelativeHumidity
}

type AccuWeatherIndicationInfoDto struct {
	Metric   AccuWeatherValueInfoDto
	Imperial AccuWeatherValueInfoDto
}

type AccuWeatherTemperatureDto struct {
	Minimum AccuWeatherValueInfoDto
	Maximum AccuWeatherValueInfoDto
}

type AccuWeatherValueInfoDto struct {
	Value    float64
	Unit     string
	UnitType int
}

type AccuWeatherRelativeHumidity struct {
	Minimum int
	Maximum int
	Average int
}

type AccuWeatherGeoPositionRequestDto struct {
	AccuWeatherBaseRequestDto
	Latitude  float64
	Longitude float64
}

type AccuWeatherGeoPositionResponseDto struct {
	Version int
	Key     string
}

type PrecipitationType string

const (
	PrecipitationTypeEmpty = ""
	PrecipitationTypeRain  = "Rain"
	PrecipitationTypeSnow  = "Snow"
	PrecipitationTypeIce   = "Ice"
)
