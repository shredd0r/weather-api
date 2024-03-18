package repository

import (
	"encoding/json"
	redisgo "github.com/redis/go-redis/v9"
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

var (
	prefixKey = "WeatherApi:"
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

func (r *RedisRepository[ENTITY]) Get(key string) (*ENTITY, error) {
	value, err := r.client.Get(addPrefix(key))

	r.log.Debugf("returned value from redis: %s", value)

	if err == redisgo.Nil {
		return nil, ErrNotFound
	}

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

func (r *RedisRepository[ENTITY]) Save(key string, e ENTITY) error {
	byteArrayEntity, err := json.Marshal(e)

	if err != nil {
		panic(ErrMarshal)
	}

	if err := r.client.Set(addPrefix(key), string(byteArrayEntity), 0); err != nil {
		r.log.Debugf(ErrFromRedisFormat, err)
		return err
	}

	return nil
}

func addPrefix(key string) string {
	return prefixKey + key
}

func (r *RedisRepository[ENTITY]) Delete(key string) error {
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

func (r *RedisCurrentWeatherRepository) Get(key string) (*entity.CurrentWeatherEntity, error) {
	return r.RedisRepository.Get(key)
}

func (r *RedisCurrentWeatherRepository) Save(key string, currentWeatherEntity entity.CurrentWeatherEntity) error {
	return r.RedisRepository.Save(key, currentWeatherEntity)
}

func (r *RedisCurrentWeatherRepository) Delete(key string) error {
	return r.RedisRepository.Delete(key)
}

type RedisHourlyWeatherRepository struct {
	*RedisRepository[[]entity.HourlyWeatherEntity]
}

func NewRedisHourlyWeatherRepository(log *logrus.Logger, client *redis.Client, cfg config.ExpirationDuration) *RedisHourlyWeatherRepository {
	return &RedisHourlyWeatherRepository{
		RedisRepository: newRedisRepository[[]entity.HourlyWeatherEntity](log, client, cfg),
	}
}

func (r *RedisHourlyWeatherRepository) Get(key string) (*[]entity.HourlyWeatherEntity, error) {
	return r.Get(key)
}

func (r *RedisHourlyWeatherRepository) Save(key string, hourlyWeatherEntityArr []entity.HourlyWeatherEntity) error {
	return r.Save(key, hourlyWeatherEntityArr)
}

func (r *RedisHourlyWeatherRepository) Delete(key string) error {
	return r.Delete(key)
}

type RedisDailyWeatherRepository struct {
	*RedisRepository[[]entity.DailyWeatherEntity]
}

func NewRedisDailyWeatherRepository(log *logrus.Logger, client *redis.Client, cfg config.ExpirationDuration) *RedisDailyWeatherRepository {
	return &RedisDailyWeatherRepository{
		RedisRepository: newRedisRepository[[]entity.DailyWeatherEntity](log, client, cfg),
	}
}

func (r *RedisDailyWeatherRepository) Get(key string) (*[]entity.DailyWeatherEntity, error) {
	return r.Get(key)
}

func (r *RedisDailyWeatherRepository) Save(key string, hourlyWeatherEntityArr []entity.DailyWeatherEntity) error {
	return r.Save(key, hourlyWeatherEntityArr)
}

func (r *RedisDailyWeatherRepository) Delete(key string) error {
	return r.Delete(key)
}
