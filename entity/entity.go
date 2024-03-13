package entity

type CurrentWeatherEntity struct {
	EpochTime           int64
	Visibility          string
	CurrentTemperature  float32
	MinTemperature      float32
	MaxTemperature      float32
	FillLikeTemperature float32
	IconResource        string
}

type HourlyWeatherEntity struct {
	Temperature                float32
	FeelsLikeTemperature       float32
	UVIndex                    int
	EpochTime                  int64
	ProbabilityOfPrecipitation float32
	AmountOfPrecipitation      float32
	Wind                       string
}

type DailyWeatherEntity struct {
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
	IconResource       string
}
