package storage

import (
	"context"
	"errors"
	"fmt"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/shredd0r/weather-api/client/redis"
	"github.com/shredd0r/weather-api/dto"
	"github.com/shredd0r/weather-api/log"
	"github.com/shredd0r/weather-api/util"
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
	logger log.Logger
	client redis.Client
}

func NewRedisWeatherStorage(logger log.Logger, client redis.Client) WeatherStorage {
	return &RedisWeatherStorage{
		logger: logger,
		client: client,
	}
}

func (r RedisWeatherStorage) SaveCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather dto.CurrentWeather) error {
	r.logger.Debugf("try save current weather to redis cache for forecaster: %s", forecaster)
	return setObjWithInnerFieldToRedis(ctx, r.client, r.getKey(keyForCurrentWeatherFormat, forecaster), addressHash, weather)
}

func (r RedisWeatherStorage) SaveHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather []*dto.HourlyWeather) error {
	r.logger.Debugf("try save hourly weather to redis cache for forecaster: %s", forecaster)
	return setObjWithInnerFieldToRedis(ctx, r.client, r.getKey(keyForHourlyWeatherFormat, forecaster), addressHash, weather)
}

func (r RedisWeatherStorage) SaveDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, weather []*dto.DailyWeather) error {
	r.logger.Debugf("try save daily weather to redis cache for forecaster: %s", forecaster)
	return setObjWithInnerFieldToRedis(ctx, r.client, r.getKey(keyForDailyWeatherFormat, forecaster), addressHash, weather)
}

func (r RedisWeatherStorage) GetCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (*dto.CurrentWeather, error) {
	r.logger.Debugf("try get current weather from redis cache for forecaster: %s", forecaster)
	currentWeather := &dto.CurrentWeather{}
	err := getObjWithInnerFieldFromRedis[dto.CurrentWeather](ctx, r.client, r.getKey(keyForCurrentWeatherFormat, forecaster), addressHash, currentWeather)
	return currentWeather, err
}

func (r RedisWeatherStorage) GetHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (*[]*dto.HourlyWeather, error) {
	r.logger.Debugf("try get hourly weather from redis cache for forecaster: %s", forecaster)
	hourlyWeathers := &[]*dto.HourlyWeather{}
	err := getObjWithInnerFieldFromRedis[[]*dto.HourlyWeather](ctx, r.client, r.getKey(keyForHourlyWeatherFormat, forecaster), addressHash, hourlyWeathers)
	return hourlyWeathers, err
}

func (r RedisWeatherStorage) GetDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) (*[]*dto.DailyWeather, error) {
	r.logger.Debugf("try get daily weather from redis cache for forecaster: %s", forecaster)
	dailyWeathers := &[]*dto.DailyWeather{}
	err := getObjWithInnerFieldFromRedis[[]*dto.DailyWeather](ctx, r.client, r.getKey(keyForDailyWeatherFormat, forecaster), addressHash, dailyWeathers)
	return dailyWeathers, err
}

func (r RedisWeatherStorage) RemoveCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) error {
	r.logger.Debugf("try remove current weather from redis cache for forecaster: %s", forecaster)
	return r.client.HDel(ctx, r.getKey(keyForCurrentWeatherFormat, forecaster), addressHash)
}
func (r RedisWeatherStorage) RemoveHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) error {
	r.logger.Debugf("try remove hourly weather from redis cache for forecaster: %s", forecaster)
	return r.client.HDel(ctx, r.getKey(keyForHourlyWeatherFormat, forecaster), addressHash)
}
func (r RedisWeatherStorage) RemoveDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) error {
	r.logger.Debugf("try remove daily weather from redis cache for forecaster: %s", forecaster)
	return r.client.HDel(ctx, r.getKey(keyForDailyWeatherUpdatedTimeFormat, forecaster), addressHash)
}

func (r RedisWeatherStorage) SaveUpdatedTimeCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, lastTime int64) error {
	r.logger.Debugf("try save updated time for current weather from redis cache for forecaster: %s", forecaster)
	return r.client.HSet(ctx, r.getKey(keyForCurrentWeatherUpdatedTimeFormat, forecaster), addressHash, lastTime)
}

func (r RedisWeatherStorage) SaveUpdatedTimeHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, lastTime int64) error {
	r.logger.Debugf("try save updated time for hourly weather from redis cache for forecaster: %s", forecaster)
	return r.client.HSet(ctx, r.getKey(keyForHourlyWeatherUpdatedTimeFormat, forecaster), addressHash, lastTime)
}

func (r RedisWeatherStorage) SaveUpdatedTimeDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster, lastTime int64) error {
	r.logger.Debugf("try save updated time for daily weather from redis cache for forecaster: %s", forecaster)
	return r.client.HSet(ctx, r.getKey(keyForDailyWeatherUpdatedTimeFormat, forecaster), addressHash, lastTime)
}

func (r RedisWeatherStorage) GetAllLastTimeUpdatedCurrentWeather(ctx context.Context, forecaster dto.WeatherForecaster) (map[string]int64, error) {
	r.logger.Debugf("try get all last time updated for current weather from redis cache for forecaster: %s", forecaster)
	return r.getAllLastTimeUpdatedWeather(ctx, keyForCurrentWeatherUpdatedTimeFormat, forecaster)
}

func (r RedisWeatherStorage) GetAllLastTimeUpdatedHourlyWeather(ctx context.Context, forecaster dto.WeatherForecaster) (map[string]int64, error) {
	r.logger.Debugf("try get all last time updated for hourly weather from redis cache for forecaster: %s", forecaster)
	return r.getAllLastTimeUpdatedWeather(ctx, keyForHourlyWeatherUpdatedTimeFormat, forecaster)
}

func (r RedisWeatherStorage) GetAllLastTimeUpdatedDailyWeather(ctx context.Context, forecaster dto.WeatherForecaster) (map[string]int64, error) {
	r.logger.Debugf("try get all last time updated for daily weather from redis cache for forecaster: %s", forecaster)
	return r.getAllLastTimeUpdatedWeather(ctx, keyForDailyWeatherUpdatedTimeFormat, forecaster)
}

func (r RedisWeatherStorage) RemoveLastTimeUpdatedCurrentWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) error {
	r.logger.Debugf("try remove updated time for current weather from redis cache for forecaster: %s", forecaster)
	return r.client.HDel(ctx, r.getKey(keyForCurrentWeatherUpdatedTimeFormat, forecaster), addressHash)
}
func (r RedisWeatherStorage) RemoveLastTimeUpdatedHourlyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) error {
	r.logger.Debugf("try remove updated time for hourly weather from redis cache for forecaster: %s", forecaster)
	return r.client.HDel(ctx, r.getKey(keyForHourlyWeatherUpdatedTimeFormat, forecaster), addressHash)
}
func (r RedisWeatherStorage) RemoveLastTimeUpdatedDailyWeather(ctx context.Context, addressHash string, forecaster dto.WeatherForecaster) error {
	r.logger.Debugf("try remove updated time for daily weather from redis cache for forecaster: %s", forecaster)
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
