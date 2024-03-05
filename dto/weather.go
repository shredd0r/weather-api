package dto

type CurrentWeatherDto struct {
	CityName            string
	DayOfWeek           string
	Visibility          string
	Month               string
	Day                 int
	CurrentWeather      float32
	MinTemperature      float32
	MaxTemperature      float32
	FillLikeTemperature float32
	IconResource        string
	LastUpdate          int64
}

type HourlyDetailWeatherDto struct {
}

type HourlyWeatherDto struct {
	Temperature                float32
	Time                       int64
	ProbabilityOfPrecipitation float32
	IconResource               string
}

type DailyWeatherDto struct {
	AverageTemperature float32
	MinTemperature     float32
	MaxTemperature     float32
	Description        string
	Humidity           float32
	IndexUV            float32
	SunriseTime        int64
	SunsetTime         int64
	WindSpeed          float32
	RainFall           float32
	DayOfWeek          int
	IconResource       string
}
