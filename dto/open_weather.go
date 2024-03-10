package dto

type OpenWeatherRequestDto struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Units     string
	Language  string `json:"lang"`
}

type OpenWeatherForecastRequestDto struct {
	OpenWeatherRequestDto
	Count int `json:"cnt"`
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

type OpenWeatherDailyWeatherResponseDto struct {
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
	Longitude *float32 `json:"lon"`
	Latitude  *float32 `json:"lat"`
}

type OpenWeatherForecastInfoDto struct {
	Visibility                 int
	DateTimeForecast           int                              `json:"dt"`
	MainInfo                   OpenWeatherMainInfoDto           `json:"main"`
	WeatherInfos               *[]OpenWeatherWeatherInfoDto     `json:"weather"`
	CloudInfo                  OpenWeatherCloudInfoDto          `json:"clouds"`
	WindInfo                   OpenWeatherWindInfoDto           `json:"wind"`
	ProbabilityOfPrecipitation float32                          `json:"pop"`
	PrecipitationInfoForRain   *OpenWeatherPrecipitationInfoDto `json:"rain"`
	PrecipitationInfoForSnow   *OpenWeatherPrecipitationInfoDto `json:"snow"`
	SystemInfo                 OpenWeatherSystemInfoDto         `json:"sys"`
	DatetimeISO                string                           `json:"dt_txt"`
}

type OpenWeatherMainInfoDto struct {
	Temperature    float32 `json:"temp"`
	FeelsLike      float32 `json:"feels_like"`
	MinTemperature float32 `json:"temp_min"`
	MaxTemperature float32 `json:"temp_max"`
	Pressure       float32
	Humidity       int
}

type OpenWeatherPrecipitationInfoDto struct {
	VolumeLastOneHour   *float32 `json:"1h"`
	VolumeLastThreeHour *float32 `json:"3h"`
}

type OpenWeatherWeatherInfoDto struct {
	Id          int
	Main        string
	Description string
	Icon        string
}

type OpenWeatherWindInfoDto struct {
	Speed   float32
	Degrees float32 `json:"deg"`
	Gust    float32
}

type OpenWeatherSystemInfoDto struct {
	Type    int
	Id      int
	Message float32
	Country string
	Sunrise int
	Sunset  int
}
