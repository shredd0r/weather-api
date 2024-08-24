package mapper

import "github.com/shredd0r/weather-api/dto"

func AccuWeatherMapToCurrentWeatherDto(accuWeather dto.AccuWeatherCurrentResponseDto) dto.CurrentWeather {
	return dto.CurrentWeather{}
}

// Link to icon number https://developer.accuweather.com/weather-icons
func AccuWeatherGetIconResource(weatherIcon uint8) dto.IconId {
	switch {
	case isClearDay(weatherIcon):
		return dto.IconIdClearDay
	case isCloudyDay(weatherIcon):
		return dto.IconIdCloudyDay
	case weatherIcon >= 7 && weatherIcon <= 8:
		return dto.IconIdCloudy
	case weatherIcon == 12 || weatherIcon == 18:
		return dto.IconIdRainy
	case weatherIcon >= 13 && weatherIcon <= 14:
		return dto.IconIdRainyDay
	case weatherIcon >= 15 && weatherIcon <= 17 || weatherIcon >= 41 && weatherIcon <= 42:
		return dto.IconIdThunder
	case weatherIcon == 19 || weatherIcon == 22 || weatherIcon >= 24 || weatherIcon <= 29:
		return dto.IconIdSnowy
	case weatherIcon == 20 || weatherIcon == 21 || weatherIcon == 23:
		return dto.IconIdSnowyDay

	}

	return dto.IconIdUndefined
}

func isClearDay(weatherIcon uint8) bool {
	return weatherIcon == 1
}

func isCloudyDay(weatherIcon uint8) bool {
	return weatherIcon >= 2 && weatherIcon <= 6
}
