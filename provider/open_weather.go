package provider

import (
	"context"
	"weather-api/client/http"
	"weather-api/dto"
	"weather-api/log"
	"weather-api/util"
)

const dayInResponse = 5

type OpenWeatherProvider struct {
	logger log.Logger
	client http.OpenWeatherInterface
}

func NewOpenWeatherProvider(logger log.Logger, client http.OpenWeatherInterface) WeatherProvider {
	return &OpenWeatherProvider{
		logger: logger,
		client: client,
	}
}

func (p OpenWeatherProvider) CurrentWeather(ctx context.Context, request *dto.WeatherRequestProviderDto) (*dto.CurrentWeather, error) {
	resp, err := p.client.GetCurrentWeatherInfo(p.getRequest(request))

	if err != nil {
		return nil, err
	}

	return &dto.CurrentWeather{
		EpochTime:            resp.EpochTime,
		Visibility:           util.PercentToFloat64Pointer(&resp.Visibility),
		CurrentTemperature:   &resp.Main.Temperature,
		MinTemperature:       &resp.Main.MinTemperature,
		MaxTemperature:       &resp.Main.MaxTemperature,
		FeelsLikeTemperature: &resp.Main.FeelsLike,
		IconResource:         nil,
	}, nil
}

func (p OpenWeatherProvider) HourlyWeather(ctx context.Context, request *dto.WeatherRequestProviderDto) (*[]*dto.HourlyWeather, error) {
	resp, err := p.client.GetForecastWeatherInfo(p.getForecastRequest(request))

	if err != nil {
		return nil, err
	}

	return p.mapToHourlyWeathers(resp), nil
}

func (p OpenWeatherProvider) DailyWeather(ctx context.Context, request *dto.WeatherRequestProviderDto) (*[]*dto.DailyWeather, error) {
	resp, err := p.client.GetForecastWeatherInfo(p.getForecastRequest(request))

	if err != nil {
		return nil, err
	}

	return p.mapToDailyWeathers(resp), nil
}

func (p OpenWeatherProvider) getUnits(unit dto.Unit) dto.OpenWeatherUnits {
	switch unit {
	case dto.UnitMetric:
		return dto.OpenWeatherUnitsMetric
	case dto.UnitImperial:
		return dto.OpenWeatherUnitsImperial
	default:
		return dto.OpenWeatherUnitsStandard
	}
}

func (p OpenWeatherProvider) mapToHourlyWeathers(response *dto.OpenWeatherHourlyWeatherResponseDto) *[]*dto.HourlyWeather {
	lenArray := len(response.ListHourlyInfo) / dayInResponse
	hourlyWeathers := make([]*dto.HourlyWeather, lenArray)

	for index := 0; index < lenArray; index++ {
		hourlyWeathers[index] = p.mapOpenWeatherForecastInfoDtoToHourlyWeather(response.ListHourlyInfo[index])
	}

	return &hourlyWeathers
}

func (p OpenWeatherProvider) mapOpenWeatherForecastInfoDtoToHourlyWeather(weatherForecast dto.OpenWeatherForecastInfoDto) *dto.HourlyWeather {
	precipitationType := p.getPrecipitationType(weatherForecast)

	return &dto.HourlyWeather{
		Temperature:                &weatherForecast.MainInfo.Temperature,
		FeelsLikeTemperature:       &weatherForecast.MainInfo.FeelsLike,
		UVIndex:                    nil,
		EpochTime:                  weatherForecast.DateTimeForecast,
		ProbabilityOfPrecipitation: &weatherForecast.ProbabilityOfPrecipitation,
		PrecipitationType:          precipitationType,
		AmountOfPrecipitation:      p.getAmountOfPrecipitation(precipitationType, weatherForecast),
		Wind: &dto.Wind{
			Speed:   &weatherForecast.WindInfo.Speed,
			Degrees: weatherForecast.WindInfo.Degrees,
		},
		IconResource: &(*weatherForecast.WeatherInfos)[0].Icon,
	}
}

func (p OpenWeatherProvider) mapToDailyWeathers(response *dto.OpenWeatherHourlyWeatherResponseDto) *[]*dto.DailyWeather {
	lenArray := dayInResponse
	dailyWeathers := make([]*dto.DailyWeather, lenArray)

	for index := 0; index < len(response.ListHourlyInfo); index = index + 3 {
		dailyWeathers[index] = p.mapOpenWeatherForecastInfoDtoToDailyWeather(response.ListHourlyInfo[index])
	}

	return &dailyWeathers
}

func (p OpenWeatherProvider) mapOpenWeatherForecastInfoDtoToDailyWeather(weatherForecast dto.OpenWeatherForecastInfoDto) *dto.DailyWeather {
	var minTemperature *float64
	var maxTemperature *float64
	precipitationType := p.getPrecipitationType(weatherForecast)

	if weatherForecast.MainInfo.MinTemperature != 0 {
		minTemperature = &weatherForecast.MainInfo.MinTemperature
	}

	if weatherForecast.MainInfo.MaxTemperature != 0 {
		maxTemperature = &weatherForecast.MainInfo.MaxTemperature
	}

	return &dto.DailyWeather{
		EpochTime:      weatherForecast.DateTimeForecast,
		MinTemperature: minTemperature,
		MaxTemperature: maxTemperature,
		Humidity:       util.PercentToFloat64(weatherForecast.MainInfo.Humidity),
		SunriseTime:    int64(weatherForecast.SystemInfo.Sunrise),
		SunsetTime:     int64(weatherForecast.SystemInfo.Sunrise),
		Wind: &dto.Wind{
			Speed:   &weatherForecast.WindInfo.Speed,
			Degrees: weatherForecast.WindInfo.Degrees,
		},
		ProbabilityOfPrecipitation: &weatherForecast.ProbabilityOfPrecipitation,
		PrecipitationType:          precipitationType,
		IconResource:               &(*weatherForecast.WeatherInfos)[0].Icon,
	}
}
func (p OpenWeatherProvider) getPrecipitationType(weatherForecast dto.OpenWeatherForecastInfoDto) dto.PrecipitationType {
	if weatherForecast.PrecipitationInfoForRain != nil {
		return dto.PrecipitationTypeRain
	}
	if weatherForecast.PrecipitationInfoForSnow != nil {
		return dto.PrecipitationTypeSnow
	}

	return dto.PrecipitationTypeEmpty
}

func (p OpenWeatherProvider) getAmountOfPrecipitation(precipitationType dto.PrecipitationType, weatherForecast dto.OpenWeatherForecastInfoDto) *float64 {
	switch precipitationType {
	case dto.PrecipitationTypeRain:
		{
			return weatherForecast.PrecipitationInfoForRain.VolumeLastThreeHour
		}
	case dto.PrecipitationTypeSnow:
		{
			return weatherForecast.PrecipitationInfoForSnow.VolumeLastThreeHour
		}
	default:
		return nil
	}
}

func (p OpenWeatherProvider) getForecastRequest(request *dto.WeatherRequestProviderDto) dto.OpenWeatherForecastRequestDto {
	return dto.OpenWeatherForecastRequestDto{
		OpenWeatherRequestDto: p.getRequest(request),
	}
}

func (p OpenWeatherProvider) getRequest(request *dto.WeatherRequestProviderDto) dto.OpenWeatherRequestDto {
	return dto.OpenWeatherRequestDto{
		Latitude:  request.Location.Coords.Latitude,
		Longitude: request.Location.Coords.Longitude,
		Language:  request.Locale,
		Units:     p.getUnits(request.Unit),
	}
}
