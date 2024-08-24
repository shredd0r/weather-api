package mapper

import (
	"fmt"
	"time"

	"github.com/shredd0r/weather-api/dto"
	"github.com/shredd0r/weather-api/util"
)

func AccuWeatherMapRequestForClient(request *dto.WeatherRequestProvider) dto.AccuWeatherRequestDto {

	return dto.AccuWeatherRequestDto{
		AccuWeatherBaseRequestDto: dto.AccuWeatherBaseRequestDto{
			Language: request.Locale,
			Metric:   request.Unit == dto.UnitMetric,
			Details:  true,
		},
		LocationKey: request.Location.AccuWeatherLocationKey,
	}
}

func AccuWeatherMapToCurrentWeatherDto(unit dto.Unit, accuWeather *dto.AccuWeatherCurrentResponseDto) *dto.CurrentWeather {
	return &dto.CurrentWeather{
		EpochTime:            accuWeather.EpochTime,
		Visibility:           getVolumeByUnit(unit, accuWeather.Visibility),
		CurrentTemperature:   getVolumeByUnit(unit, accuWeather.Temperature),
		FeelsLikeTemperature: getVolumeByUnit(unit, accuWeather.RealFeelTemperature),
		IconId:               accuWeatherGetPointerIconId(accuWeather.WeatherIcon),
		MobileLink:           accuWeather.MobileLink,
		Link:                 accuWeather.Link,
	}
}

func AccuWeatherMapToHourlyWeatherDtos(responses []*dto.AccuWeatherHourlyResponseDto) *[]*dto.HourlyWeather {
	hourlyResps := make([]*dto.HourlyWeather, len(responses))

	for index := range responses {
		hourlyResps[index] = accuWeatherMapHourlyResponse(responses[index])
	}

	return &hourlyResps
}

func AccuWeatherMapToDailyWeatherDtos(response *dto.AccuWeatherDailyResponseDto) *[]*dto.DailyWeather {
	dailyWeathers := make([]*dto.DailyWeather, len(response.DailyForecasts))

	for index := range response.DailyForecasts {
		dailyWeathers[index] = accuWeatherMapDailyResponse(&response.DailyForecasts[index])
	}

	return &dailyWeathers
}

func accuWeatherMapHourlyResponse(responseDto *dto.AccuWeatherHourlyResponseDto) *dto.HourlyWeather {
	var probabilityOfPrecipitation *float64

	if responseDto.PrecipitationProbability != nil {
		probabilityOfPrecipitation = util.PercentToFloat64Pointer(responseDto.PrecipitationProbability)
	}

	return &dto.HourlyWeather{
		CurrentTemperature:         responseDto.Temperature.Value,
		FeelsLikeTemperature:       responseDto.RealFeelTemperature.Value,
		UVIndex:                    util.PercentToFloat64Pointer(responseDto.UVIndex),
		EpochTime:                  responseDto.EpochTime,
		ProbabilityOfPrecipitation: probabilityOfPrecipitation,
		PrecipitationType:          responseDto.PrecipitationType,
		AmountOfPrecipitation:      accuWeatherGetAmountOfPrecipitation(responseDto),
		Wind: &dto.Wind{
			Speed:   responseDto.Wind.Speed.Value,
			Degrees: responseDto.Wind.Direction.Degrees,
		},
		IconId:     accuWeatherGetPointerIconId(responseDto.WeatherIcon),
		MobileLink: responseDto.MobileLink,
		Link:       responseDto.Link,
	}
}

func accuWeatherMapDailyResponse(dailyForecast *dto.AccuWeatherDailyForecastDto) *dto.DailyWeather {
	dayInfo := dto.AccuWeatherDayInfoDto{}

	if isDay(dailyForecast.EpochDate) {
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
		IconId:         accuWeatherGetPointerIconId(&dayInfo.Icon),
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

func isDay(epochTime int64) bool {
	t := time.Unix(epochTime, 0)

	return t.Hour() >= 4 && t.Hour() <= 19
}

func accuWeatherGetAmountOfPrecipitation(responseDto *dto.AccuWeatherHourlyResponseDto) *float64 {
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

func getVolumeByUnit(unit dto.Unit, v dto.AccuWeatherIndicationInfoDto) *float64 {
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

func accuWeatherGetPointerIconId(weatherIcon *int) *dto.IconId {
	if weatherIcon != nil {
		iconId := accuWeatherGetIconId(*weatherIcon)
		return &iconId
	}
	return nil
}

// Link to icon number https://developer.accuweather.com/weather-icons
func accuWeatherGetIconId(weatherIcon int) dto.IconId {
	switch {
	case isClearDay(weatherIcon):
		return dto.IconIdClearDay
	case isCloudyDay(weatherIcon):
		return dto.IconIdCloudyDay
	case isCloudy(weatherIcon):
		return dto.IconIdCloudy
	case isRainy(weatherIcon):
		return dto.IconIdRainy
	case isRainyDay(weatherIcon):
		return dto.IconIdRainyDay
	case isTrunder(weatherIcon):
		return dto.IconIdThunder
	case isSnowy(weatherIcon):
		return dto.IconIdSnowy
	case isSnowyDay(weatherIcon):
		return dto.IconIdSnowyDay
	case isClearNight(weatherIcon):
		return dto.IconIdClearNight
	case isCloudyNight(weatherIcon):
		return dto.IconIdCloudyNight
	case isRainyNight(weatherIcon):
		return dto.IconIdRainyNight
	case isSnowyNight(weatherIcon):
		return dto.IconIdSnowyNight
	case isMist(weatherIcon):
		return dto.IconIdMist
	}

	return dto.IconIdUndefined
}

func isClearDay(weatherIcon int) bool {
	return weatherIcon == 1
}

func isCloudyDay(weatherIcon int) bool {
	return weatherIcon >= 2 && weatherIcon <= 6
}

func isCloudy(weatherIcon int) bool {
	return weatherIcon == 7 || weatherIcon == 8
}

func isRainy(weatherIcon int) bool {
	return weatherIcon == 12 || weatherIcon == 18
}

func isRainyDay(weatherIcon int) bool {
	return weatherIcon == 13 || weatherIcon == 14
}

func isTrunder(weatherIcon int) bool {
	return weatherIcon >= 15 && weatherIcon <= 17 || weatherIcon == 41 || weatherIcon == 42
}

func isSnowy(weatherIcon int) bool {
	return weatherIcon == 19 || weatherIcon == 22 || weatherIcon >= 24 || weatherIcon <= 29
}

func isSnowyDay(weatherIcon int) bool {
	return weatherIcon == 20 || weatherIcon == 21 || weatherIcon == 23
}

func isClearNight(weatherIcon int) bool {
	return weatherIcon == 33
}

func isCloudyNight(weatherIcon int) bool {
	return weatherIcon >= 34 && weatherIcon <= 38
}

func isRainyNight(weatherIcon int) bool {
	return weatherIcon == 39 || weatherIcon == 40
}

func isSnowyNight(weatherIcon int) bool {
	return weatherIcon == 43 || weatherIcon == 44
}

func isMist(weatherIcon int) bool {
	return weatherIcon == 11
}
