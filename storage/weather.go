package storage

import (
	"context"
	"fmt"
	"weather-api/client/redis"
	"weather-api/dto"
)

const (
	keyForCurrentWeatherFormat            = "weather-api:current-weather:%s"         // first forecaster, second addressHash
	keyForHourlyWeatherFormat             = "weather-api:hourly-weather:%s"          //
	keyForDailyWeatherFormat              = "weather-api:daily-weather:%s"           //
	keyForCurrentWeatherUpdatedTimeFormat = "weather-api:current-weather:updated:%s" //
	keyForHourlyWeatherUpdatedTimeFormat  = "weather-api:hourly-weather:updated:%s"  //
	keyForDailyWeatherUpdatedTimeFormat   = "weather-api:daily-weather:updated:%s"   //
)

type RedisWeatherStorage struct {
	client redis.Client
}

func NewRedisWeatherStorage(client redis.Client) WeatherStorage {
	return &RedisWeatherStorage{
		client: client,
	}
}

func (r RedisWeatherStorage) GetCurrentWeatherBy(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (*dto.CurrentWeather, error) {
	currentWeather := &dto.CurrentWeather{}
	err := getObjWithInnerFieldFromRedis[dto.CurrentWeather](ctx, r.client, r.getKey(keyForCurrentWeatherFormat, forecaster), addressHash, currentWeather)
	return currentWeather, err
}

func (r RedisWeatherStorage) GetHourlyWeatherBy(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (*[]*dto.HourlyWeather, error) {
	hourlyWeathers := &[]*dto.HourlyWeather{}
	err := getObjWithInnerFieldFromRedis[[]*dto.HourlyWeather](ctx, r.client, r.getKey(keyForHourlyWeatherFormat, forecaster), addressHash, hourlyWeathers)
	return hourlyWeathers, err
}

func (r RedisWeatherStorage) GetDailyWeatherBy(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (*[]*dto.DailyWeather, error) {
	dailyWeathers := &[]*dto.DailyWeather{}
	err := getObjWithInnerFieldFromRedis[[]*dto.DailyWeather](ctx, r.client, r.getKey(keyForDailyWeatherFormat, forecaster), addressHash, dailyWeathers)
	return dailyWeathers, err
}

func (r RedisWeatherStorage) SaveCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather dto.CurrentWeather) error {
	return setObjWithInnerFieldToRedis(ctx, r.client, r.getKey(keyForCurrentWeatherFormat, forecaster), addressHash, weather)
}

func (r RedisWeatherStorage) SaveHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather []*dto.HourlyWeather) error {
	return setObjWithInnerFieldToRedis(ctx, r.client, r.getKey(keyForHourlyWeatherFormat, forecaster), addressHash, weather)
}

func (r RedisWeatherStorage) SaveDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather []*dto.DailyWeather) error {
	return setObjWithInnerFieldToRedis(ctx, r.client, r.getKey(keyForDailyWeatherFormat, forecaster), addressHash, weather)
}

func (r RedisWeatherStorage) GetLastTimeUpdatedCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (int64, error) {
	return getInt64WithInnerFieldsFromRedis(ctx, r.client, r.getKey(keyForCurrentWeatherUpdatedTimeFormat, forecaster), addressHash)
}

func (r RedisWeatherStorage) GetLastTimeUpdatedHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (int64, error) {
	return getInt64WithInnerFieldsFromRedis(ctx, r.client, r.getKey(keyForHourlyWeatherUpdatedTimeFormat, forecaster), addressHash)
}

func (r RedisWeatherStorage) GetLastTimeUpdatedDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (int64, error) {
	return getInt64WithInnerFieldsFromRedis(ctx, r.client, r.getKey(keyForDailyWeatherUpdatedTimeFormat, forecaster), addressHash)
}

func (r RedisWeatherStorage) SaveUpdatedTimeCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, lastTime int64) error {
	return r.client.HSet(ctx, r.getKey(keyForCurrentWeatherUpdatedTimeFormat, forecaster), addressHash, lastTime)
}

func (r RedisWeatherStorage) SaveUpdatedTimeHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, lastTime int64) error {
	return r.client.HSet(ctx, r.getKey(keyForHourlyWeatherUpdatedTimeFormat, forecaster), addressHash, lastTime)
}

func (r RedisWeatherStorage) SaveUpdatedTimeDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, lastTime int64) error {
	return r.client.HSet(ctx, r.getKey(keyForDailyWeatherUpdatedTimeFormat, forecaster), addressHash, lastTime)
}

func (r RedisWeatherStorage) getKey(formatKey string, forecaster dto.WeatherForecaster) string {
	return fmt.Sprintf(formatKey, forecaster)
}
