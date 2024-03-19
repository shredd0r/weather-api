package provider

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"weather-api/client/http"
	"weather-api/dto"
	"weather-api/util"
)

type HttpAccuWeatherProvider struct {
	client http.AccuWeatherInterface
	log    *logrus.Logger
}

func NewHttpAccuWeatherProvider(client http.AccuWeatherInterface, log *logrus.Logger) *HttpAccuWeatherProvider {
	return &HttpAccuWeatherProvider{
		client: client,
		log:    log,
	}
}

func (p *HttpAccuWeatherProvider) CurrentWeatherInfo(weatherRequestDto dto.WeatherRequestDto) (*dto.CurrentWeatherDto, error) {
	resp, err := p.client.GetCurrentWeatherInfo(getRequestBy(weatherRequestDto))

	if err != nil {
		return nil, err
	}

	currentWeatherDto := mapToCurrentWeatherDtoBy(resp, weatherRequestDto.Unit)

	return &currentWeatherDto, nil
}

func (p *HttpAccuWeatherProvider) HourlyWeatherInfo(weatherRequestDto dto.WeatherRequestDto) (*[]dto.HourlyWeatherDto, error) {
	resp, err := p.client.GetHourlyWeatherInfo(getRequestBy(weatherRequestDto))

	if err != nil {
		return nil, err
	}

	var hourlyDtos = make([]dto.HourlyWeatherDto, len(*resp))

	for _, respDto := range *resp {
		hourlyDtos = append(hourlyDtos, mapToHourlyWeatherDtoBy(&respDto))
	}

	return &hourlyDtos, nil
}

func (p *HttpAccuWeatherProvider) DailyWeatherInfo(weatherRequestDto dto.WeatherRequestDto) (*[]dto.DailyWeatherDto, error) {
	resp, err := p.client.GetDailyWeatherInfo(getRequestBy(weatherRequestDto))

	if err != nil {
		return nil, err
	}

	dailyWeatherSlice := mapToDailyWeatherDtoList(resp)

	return &dailyWeatherSlice, nil
}

func getRequestBy(weatherRequestDto dto.WeatherRequestDto) dto.AccuWeatherRequestDto {
	return dto.AccuWeatherRequestDto{
		LocationKey: weatherRequestDto.AccuWeatherLocationKey,
		AccuWeatherBaseRequestDto: dto.AccuWeatherBaseRequestDto{
			Language: string(weatherRequestDto.Language),
			Details:  true,
			Metric:   weatherRequestDto.Unit == dto.UnitMetric,
		},
	}
}

func mapToCurrentWeatherDtoBy(resp *dto.AccuWeatherCurrentResponseDto, unit dto.Unit) dto.CurrentWeatherDto {
	var visibility *float64
	var currentTemp *float64
	var feelsLikeTemp *float64

	switch unit {
	case dto.UnitMetric:
		{
			visibility = resp.Visibility.Metric.Value
			currentTemp = resp.Temperature.Metric.Value
			feelsLikeTemp = resp.RealFeelTemperature.Metric.Value
		}
	case dto.UnitImperial:
		{
			visibility = resp.Visibility.Imperial.Value
			currentTemp = resp.Temperature.Imperial.Value
			feelsLikeTemp = resp.RealFeelTemperature.Imperial.Value
		}
	default:
		panic(fmt.Sprintf("unspecified unit: %s", unit))
	}

	return dto.CurrentWeatherDto{
		EpochTime:            resp.EpochTime,
		Visibility:           util.RoundFloat64(visibility, 2),
		CurrentTemperature:   util.RoundFloat64(currentTemp, 2),
		FeelsLikeTemperature: util.RoundFloat64(feelsLikeTemp, 2),
		IconResource:         util.GetIconResourceNameByAccuWeatherIcon(resp.WeatherIcon),
		MobileLink:           resp.MobileLink,
		Link:                 resp.Link,
	}
}

func mapToHourlyWeatherDtoBy(hourlyDto *dto.AccuWeatherHourlyResponseDto) dto.HourlyWeatherDto {
	var precipitationType dto.PrecipitationType = dto.PrecipitationTypeEmpty

	if hourlyDto.HasPrecipitation {
		precipitationType = hourlyDto.PrecipitationType

	}

	return dto.HourlyWeatherDto{
		Temperature:                hourlyDto.Temperature.Value,
		FeelsLikeTemperature:       hourlyDto.RealFeelTemperature.Value,
		UVIndex:                    hourlyDto.UVIndex,
		EpochTime:                  hourlyDto.EpochTime,
		ProbabilityOfPrecipitation: percentToFloat64(hourlyDto.PrecipitationProbability),
		PrecipitationType:          precipitationType,
		AmountOfPrecipitation:      hourlyDto.TotalLiquid.Value,
		IconResource:               util.GetIconResourceNameByAccuWeatherIcon(hourlyDto.WeatherIcon),
		MobileLink:                 hourlyDto.MobileLink,
		Link:                       hourlyDto.Link,
		WindDto:                    mapWindDto(hourlyDto.Wind),
	}
}

func mapToDailyWeatherDtoList(dailyResp *dto.AccuWeatherDailyResponseDto) []dto.DailyWeatherDto {
	var dailySlice []dto.DailyWeatherDto

	for _, dailyForecast := range dailyResp.DailyForecasts {
		dailyWeather := dto.DailyWeatherDto{
			EpochTime:                  dailyForecast.EpochDate,
			MinTemperature:             dailyForecast.Temperature.Minimum.Value,
			MaxTemperature:             dailyForecast.Temperature.Maximum.Value,
			Humidity:                   *percentToFloat64(&dailyForecast.Day.RelativeHumidity.Average),
			UVIndex:                    percentToFloat64(getUVIndex(dailyForecast.AirAndPollen).Value),
			SunriseTime:                dailyForecast.Sun.EpochRise,
			SunsetTime:                 dailyForecast.Sun.EpochSet,
			ProbabilityOfPrecipitation: percentToFloat64(dailyForecast.Day.PrecipitationProbability),
			PrecipitationType:          dailyForecast.Day.PrecipitationType,
			IconResource:               util.GetIconResourceNameByWeatherIcon(dailyForecast.Day.Icon),
			MobileLink:                 dailyForecast.MobileLink,
			Link:                       dailyForecast.Link,
			WindDto:                    mapWindDto(dailyForecast.Day.Wind),
		}

		dailySlice = append(dailySlice, dailyWeather)
	}

	return dailySlice
}

func getUVIndex(airAndPolledSlice []dto.AccuWeatherCategoryInfoDto) dto.AccuWeatherCategoryInfoDto {
	for _, airAndPolled := range airAndPolledSlice {
		if airAndPolled.Name == "UVIndex" {
			return airAndPolled
		}
	}
	return dto.AccuWeatherCategoryInfoDto{}
}

func mapWindDto(wind dto.AccuWeatherWindInfoDto) dto.WindDto {
	return dto.WindDto{
		Speed:   wind.Speed.Value,
		Degrees: wind.Direction.Degrees,
	}
}

func percentToFloat64(percent *int) *float64 {
	if percent == nil {
		return nil
	}

	f := float64(*percent / 100.0)

	return &f
}
