package provider

import (
	"context"
	"fmt"
	"time"
	"weather-api/client/http"
	"weather-api/dto"
	"weather-api/log"
	"weather-api/util"
)

type AccuWeatherProvider struct {
	logger log.Logger
	client http.AccuWeatherInterface
}

func NewAccuWeatherProvider(logger log.Logger, client http.AccuWeatherInterface) WeatherProvider {
	return &AccuWeatherProvider{
		logger: logger,
		client: client,
	}
}

func (p AccuWeatherProvider) CurrentWeather(ctx context.Context, request *dto.WeatherRequestProviderDto) (*dto.CurrentWeather, error) {
	resp, err := p.client.GetCurrentWeatherInfo(p.mapToAccuWeatherRequest(request))

	if err != nil {
		return nil, err
	}

	return &dto.CurrentWeather{
		EpochTime:            resp.EpochTime,
		Visibility:           p.getVolumeByUnit(request.Unit, resp.Visibility),
		CurrentTemperature:   p.getVolumeByUnit(request.Unit, resp.Temperature),
		FeelsLikeTemperature: p.getVolumeByUnit(request.Unit, resp.RealFeelTemperature),
		MobileLink:           resp.MobileLink,
		Link:                 resp.Link,
	}, nil
}

func (p AccuWeatherProvider) HourlyWeather(ctx context.Context, request *dto.WeatherRequestProviderDto) (*[]*dto.HourlyWeather, error) {
	resp, err := p.client.GetHourlyWeatherInfo(p.mapToAccuWeatherRequest(request))

	if err != nil {
		return nil, err
	}

	return p.mapHourlyResponses(*resp), nil
}

func (p AccuWeatherProvider) DailyWeather(ctx context.Context, request *dto.WeatherRequestProviderDto) (*[]*dto.DailyWeather, error) {
	resp, err := p.client.GetDailyWeatherInfo(p.mapToAccuWeatherRequest(request))

	if err != nil {
		return nil, err
	}

	return p.mapDailyResponse(resp), err
}

func (p AccuWeatherProvider) mapToAccuWeatherRequest(request *dto.WeatherRequestProviderDto) dto.AccuWeatherRequestDto {

	return dto.AccuWeatherRequestDto{
		AccuWeatherBaseRequestDto: dto.AccuWeatherBaseRequestDto{
			Language: request.Locale,
			Metric:   request.Unit == dto.UnitMetric,
			Details:  true,
		},
		LocationKey: request.Location.AccuWeatherLocationKey,
	}
}

func (p AccuWeatherProvider) mapHourlyResponse(responseDto *dto.AccuWeatherHourlyResponseDto) *dto.HourlyWeather {
	var probabilityOfPrecipitation *float64

	if responseDto.PrecipitationProbability != nil {
		probabilityOfPrecipitation = util.PercentToFloat64Pointer(responseDto.PrecipitationProbability)
	}

	return &dto.HourlyWeather{
		CurrentTemperature:         responseDto.Temperature.Value,
		FeelsLikeTemperature:       responseDto.RealFeelTemperature.Value,
		UVIndex:                    responseDto.UVIndex,
		EpochTime:                  responseDto.EpochTime,
		ProbabilityOfPrecipitation: probabilityOfPrecipitation,
		PrecipitationType:          responseDto.PrecipitationType,
		AmountOfPrecipitation:      p.getAmountOfPrecipitation(responseDto),
		Wind: &dto.Wind{
			Speed:   responseDto.Wind.Speed.Value,
			Degrees: responseDto.Wind.Direction.Degrees,
		},
		IconResource: nil,
		MobileLink:   responseDto.MobileLink,
		Link:         responseDto.Link,
	}
}

func (p AccuWeatherProvider) mapHourlyResponses(responses []*dto.AccuWeatherHourlyResponseDto) *[]*dto.HourlyWeather {
	hourlyResps := make([]*dto.HourlyWeather, len(responses))

	for index := range responses {
		hourlyResps[index] = p.mapHourlyResponse(responses[index])
	}

	return &hourlyResps
}

func (p AccuWeatherProvider) mapDailyResponse(response *dto.AccuWeatherDailyResponseDto) *[]*dto.DailyWeather {
	dailyWeathers := make([]*dto.DailyWeather, len(response.DailyForecasts))

	for index := range response.DailyForecasts {
		dailyWeathers[index] = p.mapDailyForecast(&response.DailyForecasts[index])
	}

	return &dailyWeathers
}

func (p AccuWeatherProvider) mapDailyForecast(dailyForecast *dto.AccuWeatherDailyForecastDto) *dto.DailyWeather {
	dayInfo := dto.AccuWeatherDayInfoDto{}

	if p.isDay(dailyForecast.EpochDate) {
		dayInfo = dailyForecast.Day
	} else {
		dayInfo = dailyForecast.Night
	}

	dailyWeather := &dto.DailyWeather{
		EpochTime:      dailyForecast.EpochDate,
		MinTemperature: dailyForecast.Temperature.Minimum.Value,
		MaxTemperature: dailyForecast.Temperature.Maximum.Value,
		Humidity:       util.PercentToFloat64(dayInfo.RelativeHumidity.Average),
		SunriseTime:    dailyForecast.Sun.EpochRise,
		SunsetTime:     dailyForecast.Sun.EpochSet,
		IconResource:   nil,
		MobileLink:     dailyForecast.MobileLink,
		Link:           dailyForecast.Link,
	}

	if dayInfo.Wind.Speed.Value != nil {
		dailyWeather.Wind = &dto.Wind{
			Speed:   dayInfo.Wind.Speed.Value,
			Degrees: dayInfo.Wind.Direction.Degrees,
		}
	}

	if dayInfo.HasPrecipitation {
		dailyWeather.ProbabilityOfPrecipitation = util.PercentToFloat64Pointer(dayInfo.PrecipitationProbability)
		dailyWeather.PrecipitationType = dayInfo.PrecipitationType
	}

	return dailyWeather
}

func (p AccuWeatherProvider) isDay(epochTime int64) bool {
	t := time.Unix(epochTime, 0)

	return t.Hour() >= 4 && t.Hour() <= 19
}

func (p AccuWeatherProvider) getVolumeByUnit(unit dto.Unit, v dto.AccuWeatherIndicationInfoDto) *float64 {
	switch unit {
	case dto.UnitMetric:
		{
			return v.Metric.Value
		}
	case dto.UnitImperial:
		{
			return v.Imperial.Value
		}
	default:
		{
			panic(fmt.Sprintf("not supported unit, %s", unit))
		}
	}
}

func (p AccuWeatherProvider) getAmountOfPrecipitation(responseDto *dto.AccuWeatherHourlyResponseDto) *float64 {
	var probability *int

	switch responseDto.PrecipitationType {
	case dto.PrecipitationTypeIce:
		{
			probability = responseDto.IceProbability
		}
	case dto.PrecipitationTypeRain:
		{
			probability = responseDto.RainProbability
		}
	case dto.PrecipitationTypeSnow:
		{
			probability = responseDto.SnowProbability
		}
	default:
		return nil
	}

	return util.PercentToFloat64Pointer(probability)
}
