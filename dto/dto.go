package dto

type CurrentWeatherDto struct {
	EpochTime            int64
	Visibility           float64
	CurrentTemperature   float64
	MinTemperature       *float64
	MaxTemperature       *float64
	FeelsLikeTemperature float64
	IconResource         *string
	MobileLink           string
	Link                 string
}

type HourlyWeatherDto struct {
	Temperature                float64
	FeelsLikeTemperature       float64
	UVIndex                    uint8
	EpochTime                  int64
	ProbabilityOfPrecipitation float64
	PrecipitationType          PrecipitationType
	AmountOfPrecipitation      float64
	WindDto                    WindDto
	IconResource               *string
	MobileLink                 string
	Link                       string
}

type DailyWeatherDto struct {
	MinTemperature float64
	MaxTemperature float64
	Humidity       float64
	IndexUV        float64
	SunriseTime    int64
	SunsetTime     int64
	WindDto        WindDto
	RainFall       float64
	DayOfWeek      int
	IconResource   *string
	MobileLink     string
	Link           string
}

type WindDto struct {
	Speed   float64
	Degrees float64
}

type LocalityDto struct {
	Latitude               float64
	Longitude              float64
	CityName               string
	AccuWeatherLocationKey string
}

type WeatherRequestDto struct {
	LocalityDto
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
