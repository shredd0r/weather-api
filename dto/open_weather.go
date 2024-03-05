package dto

type BaseOpenWeatherRequestDto struct {
	Latitude  float32
	Longitude float32
	AppKey    string
}

type OpenWeatherRequestDto struct {
	BaseOpenWeatherRequestDto
	Units    *string
	Language *string
}

type OpenWeatherCurrentWeatherResponseDto struct {
	CoordinationDto               OpenWeatherCoordinationDto `json:"coord"`
	Weather                       *[]OpenWeatherWeatherInfoDto
	Base                          string
	Main                          OpenWeatherMainInfoDto
	Visibility                    int
	Wind                          *OpenWeatherWindInfoDto
	Clouds                        *OpenWeatherCloudInfoDto
	Rain                          *OpenWeatherPrecipitationInfoDto
	Snow                          *OpenWeatherPrecipitationInfoDto
	DatetimeOfCalculationResponse int                      `json:"dt"`
	SystemInfo                    OpenWeatherSystemInfoDto `json:"sys"`
	Timezone                      int
	Id                            int
	CityName                      string `json:"city"`
	cod                           int
}

type OpenWeatherHourlyWeatherResponseDto struct {
	Cod            string
	Message        int
	Cnt            int
	ListHourlyInfo *[]OpenWeatherHourlyInfoDto  `json:"list"`
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

type OpenWeatherHourlyInfoDto struct {
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
	Temperature    *float32 `json:"temp"`
	FeelsLike      *float32 `json:"feels_like"`
	MinTemperature *float32 `json:"temp_min"`
	MaxTemperature *float32 `json:"temp_max"`
	Pressure       *float32
	Humidity       *int
}

type OpenWeatherPrecipitationInfoDto struct {
	VolumeLastOneHour   *float32 `json:"1h"`
	VolumeLastThreeHour *float32 `json:"3h"`
}

type OpenWeatherWeatherInfoDto struct {
	Id          *int
	Main        *string
	Description *string
	Icon        *string
}

type OpenWeatherWindInfoDto struct {
	Speed   *float32
	Degrees *float32 `json:"deg"`
	gust    *float32
}

type OpenWeatherSystemInfoDto struct {
	Type    *int
	Id      *int
	Message *float32
	Country *string
	Sunrise *int
	Sunset  *int
}
