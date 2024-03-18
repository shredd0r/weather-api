package provider

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"weather-api/client/http"
	"weather-api/dto"
	"weather-api/util"
)

type AccuWeatherProvider struct {
	client *http.AccuWeatherClient
	log    *logrus.Logger
}

func NewAccuWeatherProvider(client *http.AccuWeatherClient, log *logrus.Logger) *AccuWeatherProvider {
	return &AccuWeatherProvider{
		client: client,
		log:    log,
	}
}

func (p *AccuWeatherProvider) CurrentWeatherInfo(weatherRequestDto dto.WeatherRequestDto) (*dto.CurrentWeatherDto, error) {
	resp, err := p.client.GetCurrentWeatherInfo(p.getRequestBy(weatherRequestDto))

	if err != nil {
		return nil, err
	}

	currentWeatherDto := p.mapToCurrentWeatherDtoBy(resp, weatherRequestDto.Unit)

	return &currentWeatherDto, nil
}

func (p *AccuWeatherProvider) HourlyWeatherInfo(weatherRequestDto dto.WeatherRequestDto) (*[]dto.HourlyWeatherDto, error) {
	resp, err := p.client.GetHourlyWeatherInfo(p.getRequestBy(weatherRequestDto))

	if err != nil {
		return nil, err
	}

	var hourlyDtos = make([]dto.HourlyWeatherDto, len(*resp))

	for _, respDto := range *resp {
		hourlyDtos = append(hourlyDtos, p.mapToHourlyWeatherDtoBy(&respDto))
	}

	return &hourlyDtos, nil
}

func (p *AccuWeatherProvider) DailyWeatherInfo(weatherRequestDto dto.WeatherRequestDto) (*[]dto.DailyWeatherDto, error) {
	resp, err := p.client.GetDailyWeatherInfo(p.getRequestBy(weatherRequestDto))

	if err != nil {
		return nil, err
	}

	dailyWeatherSlice := p.mapToDailyWeatherDtoList(resp)

	return &dailyWeatherSlice, nil
}

func (p *AccuWeatherProvider) getRequestBy(weatherRequestDto dto.WeatherRequestDto) dto.AccuWeatherRequestDto {
	return dto.AccuWeatherRequestDto{
		LocationKey: weatherRequestDto.AccuWeatherLocationKey,
		AccuWeatherBaseRequestDto: dto.AccuWeatherBaseRequestDto{
			Language: string(weatherRequestDto.Language),
			Details:  true,
			Metric:   weatherRequestDto.Unit == dto.UnitMetric,
		},
	}
}

func (p *AccuWeatherProvider) mapToCurrentWeatherDtoBy(resp *dto.AccuWeatherCurrentResponseDto, unit dto.Unit) dto.CurrentWeatherDto {
	var visibility float64
	var currentTemp float64
	var feelsLikeTemp float64

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

func (p *AccuWeatherProvider) mapToHourlyWeatherDtoBy(hourlyDto *dto.AccuWeatherHourlyResponseDto) dto.HourlyWeatherDto {
	var precipitationType dto.PrecipitationType = dto.PrecipitationTypeEmpty

	if hourlyDto.HasPrecipitation {
		precipitationType = hourlyDto.PrecipitationType

	}

	return dto.HourlyWeatherDto{
		Temperature:                hourlyDto.Temperature.Value,
		FeelsLikeTemperature:       hourlyDto.RealFeelTemperature.Value,
		UVIndex:                    hourlyDto.UVIndex,
		EpochTime:                  hourlyDto.EpochTime,
		ProbabilityOfPrecipitation: float64(hourlyDto.PrecipitationProbability / 100.0),
		PrecipitationType:          precipitationType,
		AmountOfPrecipitation:      hourlyDto.TotalLiquid.Value,
		IconResource:               util.GetIconResourceNameByAccuWeatherIcon(hourlyDto.WeatherIcon),
		MobileLink:                 hourlyDto.MobileLink,
		Link:                       hourlyDto.Link,
		WindDto: dto.WindDto{
			Speed:   hourlyDto.Wind.Speed.Value,
			Degrees: hourlyDto.Wind.Direction.Degrees,
		},
	}
}

func (p *AccuWeatherProvider) mapToDailyWeatherDtoList(dailyResp *dto.AccuWeatherDailyResponseDto) []dto.DailyWeatherDto {
	var dailySlice []dto.DailyWeatherDto

	for _, dailyForecast := range dailyResp.DailyForecasts {
		dailyWeather := dto.DailyWeatherDto{
			MinTemperature: dailyForecast.Temperature.Minimum.Value,
			MaxTemperature: dailyForecast.Temperature.Maximum.Value,
			//Humidity:
		}

		dailySlice = append(dailySlice, dailyWeather)
	}

	return dailySlice
}
