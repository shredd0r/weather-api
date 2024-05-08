package provider

import (
	"weather-api/dto"
	"weather-api/entity"
	"weather-api/repository"
)

type LocalityProvider interface {
	GetLocalityDto(cityName string) (*dto.LocalityDto, error)
	GetCityName(latitude float64, longitude float64) (*string, error)
}

type WeatherProvider interface {
	CurrentWeatherInfo(weatherRequestDto dto.WeatherRequestDto) (*dto.CurrentWeatherDto, error)
	HourlyWeatherInfo(weatherRequestDto dto.WeatherRequestDto) ([]*dto.HourlyWeatherDto, error)
	DailyWeatherInfo(weatherRequestDto dto.WeatherRequestDto) ([]*dto.DailyWeatherDto, error)
}

type SaveInfoProvider[T any] interface {
	Save(info T) error
}

type CacheWeatherProvider struct {
	forecaster        dto.WeatherForecaster
	currentRepository repository.WeatherRepository[entity.CurrentWeatherEntity]
	hourlyRepository  repository.WeatherRepository[entity.HourlyWeatherEntity]
	dailyRepository   repository.WeatherRepository[entity.DailyWeatherEntity]
}

func (p *CacheWeatherProvider) CurrentWeatherInfo(request dto.WeatherRequestDto) (*dto.CurrentWeatherDto, error) {
	currentEntity, err := p.currentRepository.FindByCityNameAndForecaster(request.CityName, p.forecaster)
	if err != nil {
		return nil, err
	}
	return p.mapCurrentEntityToDto(currentEntity), nil
}
func (p *CacheWeatherProvider) HourlyWeatherInfo(request dto.WeatherRequestDto) ([]*dto.HourlyWeatherDto, error) {
	hourlyEntities, err := p.hourlyRepository.FindAllByCityNameAndForecaster(request.CityName, p.forecaster)
	if err != nil {
		return nil, err
	}
	return p.mapHourlyEntitiesToDtos(hourlyEntities), nil
}
func (p *CacheWeatherProvider) DailyWeatherInfo(request dto.WeatherRequestDto) ([]*dto.DailyWeatherDto, error) {
	dailyEntities, err := p.dailyRepository.FindAllByCityNameAndForecaster(request.CityName, p.forecaster)
	if err != nil {
		return nil, err
	}
	return p.mapDailyEntitiesToDtos(dailyEntities), nil
}

func (p *CacheWeatherProvider) mapCurrentEntityToDto(currentEntity *entity.CurrentWeatherEntity) *dto.CurrentWeatherDto {
	return &dto.CurrentWeatherDto{
		EpochTime:            currentEntity.EpochTime,
		Visibility:           currentEntity.Visibility,
		CurrentTemperature:   currentEntity.CurrentTemperature,
		MinTemperature:       currentEntity.MinTemperature,
		MaxTemperature:       currentEntity.MaxTemperature,
		FeelsLikeTemperature: currentEntity.FeelsLikeTemperature,
		IconResource:         currentEntity.IconResource,
		MobileLink:           currentEntity.MobileLink,
		Link:                 currentEntity.Link,
	}
}

func (p *CacheWeatherProvider) mapWindEntityToDto(windEntity *entity.WindEntity) *dto.WindDto {
	if windEntity == nil {
		return &dto.WindDto{}
	}

	return &dto.WindDto{
		Speed:   windEntity.Speed,
		Degrees: windEntity.Degrees,
	}
}

func (p *CacheWeatherProvider) mapHourlyEntityToDto(hourlyEntity *entity.HourlyWeatherEntity) *dto.HourlyWeatherDto {
	return &dto.HourlyWeatherDto{
		Temperature:                hourlyEntity.Temperature,
		FeelsLikeTemperature:       hourlyEntity.FeelsLikeTemperature,
		UVIndex:                    hourlyEntity.UVIndex,
		EpochTime:                  hourlyEntity.EpochTime,
		ProbabilityOfPrecipitation: hourlyEntity.ProbabilityOfPrecipitation,
		PrecipitationType:          hourlyEntity.PrecipitationType,
		AmountOfPrecipitation:      hourlyEntity.AmountOfPrecipitation,
		WindDto:                    p.mapWindEntityToDto(&hourlyEntity.Wind),
		IconResource:               hourlyEntity.IconResource,
		MobileLink:                 hourlyEntity.MobileLink,
		Link:                       hourlyEntity.Link,
	}
}

func (p *CacheWeatherProvider) mapHourlyEntitiesToDtos(entities []*entity.HourlyWeatherEntity) []*dto.HourlyWeatherDto {
	dtos := make([]*dto.HourlyWeatherDto, len(entities))

	for i := range entities {
		dtos[i] = p.mapHourlyEntityToDto(entities[i])
	}

	return dtos
}

func (p *CacheWeatherProvider) mapDailyEntityToDto(dailyEntity *entity.DailyWeatherEntity) *dto.DailyWeatherDto {
	return &dto.DailyWeatherDto{
		EpochTime:                  dailyEntity.EpochTime,
		MinTemperature:             dailyEntity.MinTemperature,
		MaxTemperature:             dailyEntity.MaxTemperature,
		Humidity:                   dailyEntity.Humidity,
		UVIndex:                    dailyEntity.UVIndex,
		SunriseTime:                dailyEntity.SunriseTime,
		SunsetTime:                 dailyEntity.SunsetTime,
		WindDto:                    p.mapWindEntityToDto(dailyEntity.Wind),
		ProbabilityOfPrecipitation: dailyEntity.ProbabilityOfPrecipitation,
		PrecipitationType:          dailyEntity.PrecipitationType,
		IconResource:               dailyEntity.IconResource,
		MobileLink:                 dailyEntity.MobileLink,
		Link:                       dailyEntity.Link,
	}
}

func (p *CacheWeatherProvider) mapDailyEntitiesToDtos(entities []*entity.DailyWeatherEntity) []*dto.DailyWeatherDto {
	dtos := make([]*dto.DailyWeatherDto, len(entities))

	for i := range entities {
		dtos[i] = p.mapDailyEntityToDto(entities[i])
	}

	return dtos
}
