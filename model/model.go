package model

type CurrentWeather struct {
	CityName            string
	dayOfWeek           string
	visibility          string
	month               string
	day                 int
	minTemperature      float32
	maxTemperature      float32
	fillLikeTemperature float32
	iconResource        string
}

type DailyWeather struct {
	averageTemperature float32
	minTemperature     float32
	maxTemperature     float32
	description        string
	humidity           string
	indexUV            float32
	sunriseTime        int64
	sunsetTime         int64
	windSpeed          float32
	rainFall           float32
	dayOfWeek          int
	iconResource       string
}

type HourlyWeather struct {
}

type WeatherForecaster string

const (
	WeatherForecasterUnspecified = ""
	WeatherForecasterAccuWeather = "AccuWeather"
	WeatherForecasterOpenWeather = "OpenWeather"
	WeatherForecasterWeatherApi  = "WeatherApi"
)
