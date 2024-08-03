package storage

import (
	"context"
	"errors"
	"fmt"
	redis2 "github.com/redis/go-redis/v9"
	"weather-api/client/redis"
	"weather-api/dto"
	"weather-api/util"
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

func (r RedisWeatherStorage) SaveCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather dto.CurrentWeather) error {
	return setObjWithInnerFieldToRedis(ctx, r.client, r.getKey(keyForCurrentWeatherFormat, forecaster), addressHash, weather)
}

func (r RedisWeatherStorage) SaveHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather []*dto.HourlyWeather) error {
	return setObjWithInnerFieldToRedis(ctx, r.client, r.getKey(keyForHourlyWeatherFormat, forecaster), addressHash, weather)
}

func (r RedisWeatherStorage) SaveDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather []*dto.DailyWeather) error {
	return setObjWithInnerFieldToRedis(ctx, r.client, r.getKey(keyForDailyWeatherFormat, forecaster), addressHash, weather)
}

func (r RedisWeatherStorage) GetCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (*dto.CurrentWeather, error) {
	currentWeather := &dto.CurrentWeather{}
	err := getObjWithInnerFieldFromRedis[dto.CurrentWeather](ctx, r.client, r.getKey(keyForCurrentWeatherFormat, forecaster), addressHash, currentWeather)
	return currentWeather, err
}

func (r RedisWeatherStorage) GetHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (*[]*dto.HourlyWeather, error) {
	hourlyWeathers := &[]*dto.HourlyWeather{}
	err := getObjWithInnerFieldFromRedis[[]*dto.HourlyWeather](ctx, r.client, r.getKey(keyForHourlyWeatherFormat, forecaster), addressHash, hourlyWeathers)
	return hourlyWeathers, err
}

func (r RedisWeatherStorage) GetDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (*[]*dto.DailyWeather, error) {
	dailyWeathers := &[]*dto.DailyWeather{}
	err := getObjWithInnerFieldFromRedis[[]*dto.DailyWeather](ctx, r.client, r.getKey(keyForDailyWeatherFormat, forecaster), addressHash, dailyWeathers)
	return dailyWeathers, err
}

func (r RedisWeatherStorage) RemoveCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) error {
	return r.client.HDel(ctx, r.getKey(keyForCurrentWeatherFormat, forecaster), addressHash)
}
func (r RedisWeatherStorage) RemoveHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) error {
	return r.client.HDel(ctx, r.getKey(keyForHourlyWeatherFormat, forecaster), addressHash)
}
func (r RedisWeatherStorage) RemoveDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) error {
	return r.client.HDel(ctx, r.getKey(keyForDailyWeatherUpdatedTimeFormat, forecaster), addressHash)
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

func (r RedisWeatherStorage) GetAllLastTimeUpdatedCurrentWeather(ctx context.Context, forecaster dto.WeatherForecaster) (map[string]int64, error) {
	return r.getAllLastTimeUpdatedWeather(ctx, keyForCurrentWeatherUpdatedTimeFormat, forecaster)
}

func (r RedisWeatherStorage) GetAllLastTimeUpdatedHourlyWeather(ctx context.Context, forecaster dto.WeatherForecaster) (map[string]int64, error) {
	return r.getAllLastTimeUpdatedWeather(ctx, keyForHourlyWeatherUpdatedTimeFormat, forecaster)
}

func (r RedisWeatherStorage) GetAllLastTimeUpdatedDailyWeather(ctx context.Context, forecaster dto.WeatherForecaster) (map[string]int64, error) {
	return r.getAllLastTimeUpdatedWeather(ctx, keyForDailyWeatherUpdatedTimeFormat, forecaster)
}

func (r RedisWeatherStorage) RemoveLastTimeUpdatedCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) error {
	return r.client.HDel(ctx, r.getKey(keyForCurrentWeatherUpdatedTimeFormat, forecaster), addressHash)
}
func (r RedisWeatherStorage) RemoveLastTimeUpdatedHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) error {
	return r.client.HDel(ctx, r.getKey(keyForHourlyWeatherUpdatedTimeFormat, forecaster), addressHash)
}
func (r RedisWeatherStorage) RemoveLastTimeUpdatedDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) error {
	return r.client.HDel(ctx, r.getKey(keyForDailyWeatherUpdatedTimeFormat, forecaster), addressHash)
}

func (r RedisWeatherStorage) getAllLastTimeUpdatedWeather(ctx context.Context, keyFormat string, forecaster dto.WeatherForecaster) (map[string]int64, error) {
	mapAddressHashToBytesLastTime, err := r.client.HGetAll(ctx, r.getKey(keyFormat, forecaster))
	if err != nil {
		if errors.Is(err, redis2.Nil) {
			return nil, ErrNotFound
		}

		return nil, err
	}
	mapAddressHashToLastTime := map[string]int64{}

	for addressToHash, bytes := range mapAddressHashToBytesLastTime {
		mapAddressHashToLastTime[addressToHash] = util.BytesToInt64(bytes)
	}

	return mapAddressHashToLastTime, nil
}

func (r RedisWeatherStorage) getKey(formatKey string, forecaster dto.WeatherForecaster) string {
	return fmt.Sprintf(formatKey, forecaster)
}
