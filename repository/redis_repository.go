package repository

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"weather-api/client/redis"
	"weather-api/config"
	"weather-api/entity"
)

var (
	ErrUnmarshal       = "unmarshal entity from redis response error"
	ErrMarshal         = "marshal entity error"
	ErrFromRedisFormat = "error from redisClient: %s"
)

type RedisRepository[ENTITY any] struct {
	log    *logrus.Logger
	cfg    config.ExpirationDuration
	client *redis.Client
}

func newRedisRepository[ENTITY any](log *logrus.Logger, client *redis.Client, cfg config.ExpirationDuration) *RedisRepository[ENTITY] {
	return &RedisRepository[ENTITY]{
		log:    log,
		cfg:    cfg,
		client: client,
	}
}

func (r *RedisRepository[ENTITY]) Get(cityName string) (*ENTITY, error) {
	value, err := r.client.Get(cityName)

	r.log.Debugf("returned value from redis: %s", value)

	if err != nil {
		r.log.Debugf(ErrFromRedisFormat, err)
		return nil, err
	}

	var e ENTITY

	if err := json.Unmarshal([]byte(value), &e); err != nil {
		panic(ErrUnmarshal)
	}

	return &e, nil
}

func (r *RedisRepository[ENTITY]) Save(cityName string, e ENTITY) error {
	byteArrayEntity, err := json.Marshal(e)

	if err != nil {
		panic(ErrMarshal)
	}

	if err := r.client.Set(cityName, string(byteArrayEntity), r.cfg.WeatherInfo); err != nil {
		r.log.Debugf(ErrFromRedisFormat, err)
		return err
	}

	return nil
}

func (r *RedisRepository[ENTITY]) Delete(cityName string) error {
	return nil
}

type RedisCurrentWeatherRepository struct {
	*RedisRepository[entity.CurrentWeatherEntity]
}

func NewRedisCurrentWeatherRepository(log *logrus.Logger, client *redis.Client, cfg config.ExpirationDuration) *RedisCurrentWeatherRepository {
	return &RedisCurrentWeatherRepository{
		RedisRepository: newRedisRepository[entity.CurrentWeatherEntity](log, client, cfg),
	}
}

func (r *RedisCurrentWeatherRepository) Get(cityName string) (*entity.CurrentWeatherEntity, error) {
	return r.Get(cityName)
}

func (r *RedisCurrentWeatherRepository) Save(cityName string, currentWeatherEntity entity.CurrentWeatherEntity) error {
	return r.Save(cityName, currentWeatherEntity)
}

func (r *RedisCurrentWeatherRepository) Delete(cityName string) error {
	return r.Delete(cityName)
}

type RedisHourlyWeatherRepository struct {
	*RedisRepository[[]entity.HourlyWeatherEntity]
}

func NewRedisHourlyWeatherRepository(log *logrus.Logger, client *redis.Client, cfg config.ExpirationDuration) *RedisHourlyWeatherRepository {
	return &RedisHourlyWeatherRepository{
		RedisRepository: newRedisRepository[[]entity.HourlyWeatherEntity](log, client, cfg),
	}
}

func (r *RedisHourlyWeatherRepository) Get(cityName string) (*[]entity.HourlyWeatherEntity, error) {
	return r.Get(cityName)
}

func (r *RedisHourlyWeatherRepository) Save(cityName string, hourlyWeatherEntityArr []entity.HourlyWeatherEntity) error {
	return r.Save(cityName, hourlyWeatherEntityArr)
}

func (r *RedisHourlyWeatherRepository) Delete(cityName string) error {
	return r.Delete(cityName)
}

type RedisDailyWeatherRepository struct {
	*RedisRepository[[]entity.DailyWeatherEntity]
}

func NewRedisDailyWeatherRepository(log *logrus.Logger, client *redis.Client, cfg config.ExpirationDuration) *RedisDailyWeatherRepository {
	return &RedisDailyWeatherRepository{
		RedisRepository: newRedisRepository[[]entity.DailyWeatherEntity](log, client, cfg),
	}
}

func (r *RedisDailyWeatherRepository) Get(cityName string) (*[]entity.DailyWeatherEntity, error) {
	return r.Get(cityName)
}

func (r *RedisDailyWeatherRepository) Save(cityName string, hourlyWeatherEntityArr []entity.DailyWeatherEntity) error {
	return r.Save(cityName, hourlyWeatherEntityArr)
}

func (r *RedisDailyWeatherRepository) Delete(cityName string) error {
	return r.Delete(cityName)
}
