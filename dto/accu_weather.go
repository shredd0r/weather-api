package dto

type AccuWeatherRequestDto struct {
	AppKey      string
	Language    string
	Metric      bool
	Details     *bool
	LocationKey string
}

type AccuWeatherHourlyResponseDto struct {
	DateTime                 string
	EpochDateTime            int64
	WeatherIcon              string
	IconPhrase               string
	HasPrecipitation         bool
	IsDayLight               bool
	Temperature              AccuWeatherTemperatureInfoDto
	PrecipitationProbability int
	MobileLink               string
	Link                     string
}

type AccuWeatherDailyResponseDto struct {
	Headline       AccuWeatherHeadlineDto
	DailyForecasts []AccuWeatherDailyForecastDto
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
	Temperature AccuWeatherTemperatureInfoDto
	Day         AccuWeatherDayInfoDto
	Night       AccuWeatherDayInfoDto
	Sources     *[]string
	MobileLink  string
	Link        string
}

type AccuWeatherDayInfoDto struct {
	Rise      string
	EpochRise int
	Set       string
	EpochSet  int
}

type AccuWeatherTemperatureInfoDto struct {
	Value    *float32
	Unit     string
	UnitType int
}

type AccuWeatherGeoPositionRequestDto struct {
	AppKey    string
	Latitude  float32
	Longitude float32
	Language  string
	Metric    *bool
	details   *bool
}

type AccuWeatherGeoPositionResponseDto struct {
	Key string
}
