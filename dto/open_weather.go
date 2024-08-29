package dto

type OpenWeatherUnits string

const (
	OpenWeatherUnitsUndefined = "undefined"
	OpenWeatherUnitsMetric    = "metric"
	OpenWeatherUnitsImperial  = "imperial"
	OpenWeatherUnitsStandard  = "standard"
)

type OpenWeatherIcon string

const (
	OpenWeatherClearSkyDay          = "01d"
	OpenWeatherClearSkyNight        = "01n"
	OpenWeatherFewCloudsDay         = "02d"
	OpenWeatherFewCloudsNight       = "02n"
	OpenWeatherScatteredCloudsDay   = "03d"
	OpenWeatherScatteredCloudsNight = "03n"
	OpenWeatherBrokenCloudsDay      = "04d"
	OpenWeatherBrokenCloudsNight    = "04n"
	OpenWeatherShowerRainDay        = "09d"
	OpenWeatherShowerRainNight      = "09n"
	OpenWeatherRainDay              = "10d"
	OpenWeatherRainNight            = "10n"
	OpenWeatherThunderstormDay      = "11d"
	OpenWeatherThunderstormNight    = "11n"
	OpenWeatherSnowDay              = "13d"
	OpenWeatherSnowNight            = "13n"
	OpenWeatherMistDay              = "50d"
	OpenWeatherMistNight            = "50n"
)

type OpenWeatherRequestDto struct {
	Latitude  float64          `json:"lat"`
	Longitude float64          `json:"lon"`
	Units     OpenWeatherUnits `json:"units"`
	Language  string           `json:"lang"`
}

type OpenWeatherForecastRequestDto struct {
	OpenWeatherRequestDto
	Count *int `json:"cnt"`
}

type OpenWeatherCurrentWeatherResponseDto struct {
	CoordinationDto OpenWeatherCoordinationDto `json:"coord"`
	Weather         []OpenWeatherWeatherInfoDto
	Base            string
	Main            OpenWeatherMainInfoDto
	Visibility      int
	Wind            OpenWeatherWindInfoDto
	Clouds          OpenWeatherCloudInfoDto
	Rain            OpenWeatherPrecipitationInfoDto
	Snow            OpenWeatherPrecipitationInfoDto
	EpochTime       int64                    `json:"dt"`
	SystemInfo      OpenWeatherSystemInfoDto `json:"sys"`
	Timezone        int64
}

type OpenWeatherHourlyWeatherResponseDto struct {
	Cod            string
	Cnt            int
	ListHourlyInfo []OpenWeatherForecastInfoDto `json:"list"`
	CityInfo       OpenWeatherCityInfoDetailDto `json:"city"`
}

type OpenWeatherCityInfoDetailDto struct {
	Id           int
	Name         string
	Country      string
	Population   int64
	Timezone     int64
	Sunrise      int64
	Sunset       int64
	Coordination OpenWeatherCoordinationDto `json:"coord"`
}

type OpenWeatherCityInfoDto struct {
	Id           int
	Name         string
	Coordination OpenWeatherCoordinationDto `json:"coord"`
}

type OpenWeatherCloudInfoDto struct {
	All *int
}

type OpenWeatherCoordinationDto struct {
	Longitude *float64 `json:"lon"`
	Latitude  *float64 `json:"lat"`
}

type OpenWeatherForecastInfoDto struct {
	Visibility                 int
	DateTimeForecast           int64                            `json:"dt"`
	MainInfo                   OpenWeatherMainInfoDto           `json:"main"`
	WeatherInfos               *[]OpenWeatherWeatherInfoDto     `json:"weather"`
	CloudInfo                  OpenWeatherCloudInfoDto          `json:"clouds"`
	WindInfo                   OpenWeatherWindInfoDto           `json:"wind"`
	ProbabilityOfPrecipitation float64                          `json:"pop"`
	PrecipitationInfoForRain   *OpenWeatherPrecipitationInfoDto `json:"rain"`
	PrecipitationInfoForSnow   *OpenWeatherPrecipitationInfoDto `json:"snow"`
	SystemInfo                 OpenWeatherSystemInfoDto         `json:"sys"`
	DatetimeISO                string                           `json:"dt_txt"`
}

type OpenWeatherMainInfoDto struct {
	Temperature    float64 `json:"temp"`
	FeelsLike      float64 `json:"feels_like"`
	MinTemperature float64 `json:"temp_min"`
	MaxTemperature float64 `json:"temp_max"`
	Pressure       float64
	Humidity       int
}

type OpenWeatherPrecipitationInfoDto struct {
	VolumeLastOneHour   *float64 `json:"1h"`
	VolumeLastThreeHour *float64 `json:"3h"`
}

type OpenWeatherWeatherInfoDto struct {
	Id          int
	Main        string
	Description string
	Icon        string
}

type OpenWeatherWindInfoDto struct {
	Speed   float64
	Degrees float64 `json:"deg"`
	Gust    float64
}

type OpenWeatherSystemInfoDto struct {
	Type    int
	Id      int
	Message float64
	Country string
	Sunrise int
	Sunset  int
}

type OpenWeatherError struct {
	Cod     int
	Message string
}
