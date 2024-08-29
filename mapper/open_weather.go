package mapper

import (
	"github.com/shredd0r/weather-api/dto"
	"github.com/shredd0r/weather-api/util"
)

const countDays = 5
const countWeatherInfoPerDay = 8

func OpenWeatherMapCurrentWeather(response *dto.OpenWeatherCurrentWeatherResponseDto) *dto.CurrentWeather {
	return &dto.CurrentWeather{
		EpochTime:            response.EpochTime,
		Visibility:           util.PercentToFloat64Pointer(&response.Visibility),
		CurrentTemperature:   &response.Main.Temperature,
		FeelsLikeTemperature: &response.Main.FeelsLike,
		WeatherId:            openWeatherGetPointerWeatherId(&response.Weather),
	}
}

func OpenWeatherMapHourlyWeahters(response *dto.OpenWeatherHourlyWeatherResponseDto) *[]*dto.HourlyWeather {
	lenArray := countWeatherInfoPerDay
	hourlyWeathers := make([]*dto.HourlyWeather, lenArray)

	for index := 0; index < lenArray; index++ {
		hourlyWeathers[index] = openWeatherForecastInfoDtoToHourlyWeather(response.ListHourlyInfo[index])
	}

	return &hourlyWeathers
}

func OpenWeatherMapDailyWeahters(response *dto.OpenWeatherHourlyWeatherResponseDto) *[]*dto.DailyWeather {
	lenArray := countDays
	dailyWeathers := make([]*dto.DailyWeather, lenArray)

	for index := 0; index < lenArray; index++ {
		dailyWeathers[index] = mapOpenWeatherForecastInfoDtoToDailyWeather(response.ListHourlyInfo[index+countWeatherInfoPerDay])
	}

	return &dailyWeathers
}

func OpenWeatherGetRequest(request *dto.WeatherRequestProvider) dto.OpenWeatherRequestDto {
	return dto.OpenWeatherRequestDto{
		Latitude:  request.Location.Coords.Latitude,
		Longitude: request.Location.Coords.Longitude,
		Language:  request.Locale,
		Units:     getUnits(request.Unit),
	}
}

func OpenWeatherGetForecastRequest(request *dto.WeatherRequestProvider) dto.OpenWeatherForecastRequestDto {
	return dto.OpenWeatherForecastRequestDto{
		OpenWeatherRequestDto: OpenWeatherGetRequest(request),
	}
}

func mapOpenWeatherForecastInfoDtoToDailyWeather(weatherForecast dto.OpenWeatherForecastInfoDto) *dto.DailyWeather {
	var minTemperature *float64
	var maxTemperature *float64
	precipitationType := openWeatherGetPrecipitationType(weatherForecast)

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
		WeatherId:                  openWeatherGetPointerWeatherId(weatherForecast.WeatherInfos),
	}
}

func openWeatherForecastInfoDtoToHourlyWeather(weatherForecast dto.OpenWeatherForecastInfoDto) *dto.HourlyWeather {
	precipitationType := openWeatherGetPrecipitationType(weatherForecast)

	return &dto.HourlyWeather{
		CurrentTemperature:         &weatherForecast.MainInfo.Temperature,
		FeelsLikeTemperature:       &weatherForecast.MainInfo.FeelsLike,
		UVIndex:                    nil,
		EpochTime:                  weatherForecast.DateTimeForecast,
		ProbabilityOfPrecipitation: &weatherForecast.ProbabilityOfPrecipitation,
		PrecipitationType:          precipitationType,
		AmountOfPrecipitation:      openWeatherGetAmountOfPrecipitation(precipitationType, weatherForecast),
		Wind: &dto.Wind{
			Speed:   &weatherForecast.WindInfo.Speed,
			Degrees: weatherForecast.WindInfo.Degrees,
		},
		WeatherId: openWeatherGetPointerWeatherId(weatherForecast.WeatherInfos),
	}
}

func openWeatherGetAmountOfPrecipitation(precipitationType dto.PrecipitationType, weatherForecast dto.OpenWeatherForecastInfoDto) *float64 {
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

func openWeatherGetPrecipitationType(weatherForecast dto.OpenWeatherForecastInfoDto) dto.PrecipitationType {
	if weatherForecast.PrecipitationInfoForRain != nil {
		return dto.PrecipitationTypeRain
	}
	if weatherForecast.PrecipitationInfoForSnow != nil {
		return dto.PrecipitationTypeSnow
	}

	return dto.PrecipitationTypeEmpty
}

func openWeatherGetPointerWeatherId(weatherInfoDtos *[]dto.OpenWeatherWeatherInfoDto) *dto.WeatherId {
	if weatherInfoDtos == nil {
		return nil
	}

	var WeatherId dto.WeatherId
	weatherInfo := (*weatherInfoDtos)[0]

	switch weatherInfo.Icon {
	case dto.OpenWeatherClearSkyDay:
		WeatherId = dto.WeatherIdClearDay
	case dto.OpenWeatherClearSkyNight:
		WeatherId = dto.WeatherIdClearNight
	case dto.OpenWeatherFewCloudsDay:
		WeatherId = dto.WeatherIdCloudyDay
	case dto.OpenWeatherFewCloudsNight:
		WeatherId = dto.WeatherIdCloudyNight
	case dto.OpenWeatherScatteredCloudsDay, dto.OpenWeatherScatteredCloudsNight, dto.OpenWeatherBrokenCloudsDay, dto.OpenWeatherBrokenCloudsNight:
		WeatherId = dto.WeatherIdCloudy
	case dto.OpenWeatherShowerRainDay, dto.OpenWeatherRainDay:
		WeatherId = dto.WeatherIdRainyDay
	case dto.OpenWeatherShowerRainNight, dto.OpenWeatherRainNight:
		WeatherId = dto.WeatherIdRainyNight
	case dto.OpenWeatherThunderstormDay, dto.OpenWeatherThunderstormNight:
		WeatherId = dto.WeatherIdThunder
	case dto.OpenWeatherSnowDay, dto.OpenWeatherSnowNight:
		WeatherId = dto.WeatherIdSnowy
	case dto.OpenWeatherMistDay, dto.OpenWeatherMistNight:
		WeatherId = dto.WeatherIdMist
	}
	return &WeatherId
}

func getUnits(unit dto.Unit) dto.OpenWeatherUnits {
	switch unit {
	case dto.UnitMetric:
		return dto.OpenWeatherUnitsMetric
	case dto.UnitImperial:
		return dto.OpenWeatherUnitsImperial
	default:
		return dto.OpenWeatherUnitsStandard
	}
}
