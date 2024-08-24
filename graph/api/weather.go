package api

import (
	"context"
	"errors"

	"github.com/shredd0r/weather-api/dto"
	"github.com/shredd0r/weather-api/graph/model"
	"github.com/shredd0r/weather-api/service"
)

var (
	ErrInvalidWeatherForecaster = errors.New("invalid weather forecaster")
)

type WeatherGraphqlApi struct {
	locationService    service.LocationService
	accuWeatherService service.WeatherService
	openWeatherService service.WeatherService
}

func NewWeatherGraphqlApi(
	locationService service.LocationService,
	accuWeatherService service.WeatherService,
	openWeatherService service.WeatherService) *WeatherGraphqlApi {
	return &WeatherGraphqlApi{
		locationService:    locationService,
		accuWeatherService: accuWeatherService,
		openWeatherService: openWeatherService,
	}
}

func (p *WeatherGraphqlApi) FindGeocoding(ctx context.Context, input *model.GeocodingRequest) (*[]*model.Geocoding, error) {
	resp, err := p.locationService.FindGeocoding(ctx, p.mapToGeocodingRequest(input))

	if err != nil {
		return nil, err
	}

	return p.mapToGeocodings(*resp), nil
}

func (p *WeatherGraphqlApi) CurrentWeather(ctx context.Context, input *model.WeatherRequest) (*model.CurrentWeather, error) {
	request := p.mapToWeatherRequest(input)
	var funcCurrentWeather func(context.Context, *dto.WeatherRequest) (*dto.CurrentWeather, error)
	switch input.Forecaster {
	case dto.WeatherForecasterAccuWeather:
		{
			funcCurrentWeather = p.accuWeatherService.CurrentWeather
		}
	case dto.WeatherForecasterOpenWeather:
		{
			funcCurrentWeather = p.openWeatherService.CurrentWeather
		}
	default:
		{
			return nil, ErrInvalidWeatherForecaster
		}
	}

	resp, err := funcCurrentWeather(ctx, request)

	if err != nil {
		return nil, err
	}

	return p.mapToCurrentWeather(resp), nil
}

func (p *WeatherGraphqlApi) HourlyWeather(ctx context.Context, input *model.WeatherRequest) (*[]*model.HourlyWeather, error) {
	request := p.mapToWeatherRequest(input)
	var funcHourlyWeather func(context.Context, *dto.WeatherRequest) (*[]*dto.HourlyWeather, error)
	switch input.Forecaster {
	case dto.WeatherForecasterAccuWeather:
		{
			funcHourlyWeather = p.accuWeatherService.HourlyWeather
		}
	case dto.WeatherForecasterOpenWeather:
		{
			funcHourlyWeather = p.openWeatherService.HourlyWeather
		}
	default:
		{
			return nil, ErrInvalidWeatherForecaster
		}
	}

	resp, err := funcHourlyWeather(ctx, request)

	if err != nil {
		return nil, err
	}

	return p.mapToHourlyWeathers(*resp), nil
}

func (p *WeatherGraphqlApi) DailyWeather(ctx context.Context, input *model.WeatherRequest) (*[]*model.DailyWeather, error) {
	request := p.mapToWeatherRequest(input)
	var funcDailyWeather func(context.Context, *dto.WeatherRequest) (*[]*dto.DailyWeather, error)
	switch input.Forecaster {
	case dto.WeatherForecasterAccuWeather:
		{
			funcDailyWeather = p.accuWeatherService.DailyWeather
		}
	case dto.WeatherForecasterOpenWeather:
		{
			funcDailyWeather = p.openWeatherService.DailyWeather
		}
	default:
		{
			return nil, ErrInvalidWeatherForecaster
		}
	}

	resp, err := funcDailyWeather(ctx, request)

	if err != nil {
		return nil, err
	}

	return p.mapToDailyWeathers(*resp), nil
}

func (p *WeatherGraphqlApi) mapToWeatherRequest(input *model.WeatherRequest) *dto.WeatherRequest {
	return &dto.WeatherRequest{
		Coords: &dto.Coords{
			Latitude:  input.Coords.Latitude,
			Longitude: input.Coords.Longitude,
		},
		Locale: input.Locale,
		Unit:   input.Unit,
	}
}

func (p *WeatherGraphqlApi) mapToGeocodingRequest(input *model.GeocodingRequest) *dto.GeocodingRequest {
	r := &dto.GeocodingRequest{
		City:    input.City,
		State:   input.State,
		Country: input.Country,
	}
	return r
}

func (p *WeatherGraphqlApi) mapToGeocoding(resp *dto.Geocoding) *model.Geocoding {
	return &model.Geocoding{
		Name:      resp.Name,
		Latitude:  resp.Latitude,
		Longitude: resp.Longitude,
		Country:   resp.Country,
		State:     resp.State,
	}
}

func (p *WeatherGraphqlApi) mapToGeocodings(resp []*dto.Geocoding) *[]*model.Geocoding {
	geocodings := make([]*model.Geocoding, len(resp))

	for index := range resp {
		geocodings[index] = p.mapToGeocoding(resp[index])
	}

	return &geocodings
}

func (p *WeatherGraphqlApi) mapToCurrentWeather(resp *dto.CurrentWeather) *model.CurrentWeather {
	return &model.CurrentWeather{
		EpochTime:            resp.EpochTime,
		Visibility:           resp.Visibility,
		CurrentTemperature:   resp.CurrentTemperature,
		MinTemperature:       resp.MinTemperature,
		MaxTemperature:       resp.MaxTemperature,
		FeelsLikeTemperature: resp.FeelsLikeTemperature,
		IconID:               resp.IconId,
		MobileLink:           resp.MobileLink,
		Link:                 resp.Link,
	}
}

func (p *WeatherGraphqlApi) mapToHourlyWeather(resp *dto.HourlyWeather) *model.HourlyWeather {
	return &model.HourlyWeather{
		EpochTime:                  resp.EpochTime,
		CurrentTemperature:         resp.CurrentTemperature,
		FeelsLikeTemperature:       resp.FeelsLikeTemperature,
		ProbabilityOfPrecipitation: resp.ProbabilityOfPrecipitation,
		PrecipitationType:          resp.PrecipitationType,
		AmountOfPrecipitation:      resp.AmountOfPrecipitation,
		UvIndex:                    resp.UVIndex,
		Wind:                       p.mapToWind(resp.Wind),
		IconID:                     resp.IconId,
		MobileLink:                 resp.MobileLink,
		Link:                       resp.Link,
	}
}

func (p *WeatherGraphqlApi) mapToHourlyWeathers(resp []*dto.HourlyWeather) *[]*model.HourlyWeather {
	modelHourlyWeather := make([]*model.HourlyWeather, len(resp))

	for index := range resp {
		modelHourlyWeather[index] = p.mapToHourlyWeather(resp[index])
	}

	return &modelHourlyWeather
}

func (p *WeatherGraphqlApi) mapToDailyWeather(resp *dto.DailyWeather) *model.DailyWeather {
	return &model.DailyWeather{
		EpochTime:                  resp.EpochTime,
		MinTemperature:             resp.MinTemperature,
		MaxTemperature:             resp.MaxTemperature,
		Humidity:                   resp.Humidity,
		UvIndex:                    resp.UVIndex,
		SunriseTime:                resp.SunriseTime,
		SunsetTime:                 resp.SunsetTime,
		Wind:                       p.mapToWind(resp.Wind),
		ProbabilityOfPrecipitation: resp.ProbabilityOfPrecipitation,
		PrecipitationType:          resp.PrecipitationType,
		IconID:                     resp.IconId,
		MobileLink:                 resp.MobileLink,
		Link:                       resp.Link,
	}
}

func (p *WeatherGraphqlApi) mapToDailyWeathers(resp []*dto.DailyWeather) *[]*model.DailyWeather {
	modelDailyWeather := make([]*model.DailyWeather, len(resp))

	for index := range resp {
		modelDailyWeather[index] = p.mapToDailyWeather(resp[index])
	}

	return &modelDailyWeather
}

func (p *WeatherGraphqlApi) mapToWind(wind *dto.Wind) *model.Wind {
	if wind != nil {
		return &model.Wind{
			Speed:   wind.Speed,
			Degrees: wind.Degrees,
		}
	}
	return nil
}
