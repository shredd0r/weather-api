package dto

type CurrentWeather struct {
	EpochTime            int64
	Visibility           *float64
	CurrentTemperature   *float64
	MinTemperature       *float64
	MaxTemperature       *float64
	FeelsLikeTemperature *float64
	IconResource         *string
	MobileLink           string
	Link                 string
}

type HourlyWeather struct {
	Temperature                *float64
	FeelsLikeTemperature       *float64
	UVIndex                    *uint8
	EpochTime                  int64
	ProbabilityOfPrecipitation *float64
	PrecipitationType          PrecipitationType
	AmountOfPrecipitation      *float64
	WindDto                    *Wind
	IconResource               *string
	MobileLink                 string
	Link                       string
}

type DailyWeather struct {
	EpochTime                  int64
	MinTemperature             *float64
	MaxTemperature             *float64
	Humidity                   float64
	UVIndex                    *float64
	SunriseTime                int64
	SunsetTime                 int64
	WindDto                    *Wind
	ProbabilityOfPrecipitation *float64
	PrecipitationType          PrecipitationType
	IconResource               *string
	MobileLink                 string
	Link                       string
}

type Wind struct {
	Speed   *float64
	Degrees float64
}

type Coords struct {
	Latitude  float64
	Longitude float64
}

type Location struct {
	Coords                 Coords
	AccuWeatherLocationKey string
}

type LocationInfo struct {
	Coords                 Coords
	AddressHash            string
	AccuWeatherLocationKey string
}

type WeatherRequestDto struct {
	Coords *Coords
	Language
	Unit
}

type WeatherRequestProviderDto struct {
	Location LocationInfo
	Language
	Unit
}

type WeatherForecaster string

const (
	WeatherForecasterUnspecified = ""
	WeatherForecasterAccuWeather = "AccuWeather"
	WeatherForecasterOpenWeather = "OpenWeather"
	WeatherForecasterWeatherApi  = "WeatherApi"
)

type Language string

const (
	LanguageUnspecified = ""
	LanguageEnglish     = "en"
	LanguageUkrainian   = "uk"
)

type Unit string

const (
	UnitUnspecified = ""
	UnitImperial    = "imperial"
	UnitMetric      = "metric"
)
