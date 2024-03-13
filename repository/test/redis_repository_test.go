package test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"weather-api/client/redis"
	"weather-api/config"
	"weather-api/entity"
	"weather-api/logger"
	redis2 "weather-api/repository"
)

var cityName = "Kharkiv"

func TestSaveCurrentWeatherEntity(t *testing.T) {
	r := initRedisCurrentRepository(t)

	err := r.Save(cityName, entity.CurrentWeatherEntity{
		EpochTime:           time.Now().UnixMilli(),
		Visibility:          "visibility",
		CurrentTemperature:  10.24,
		MinTemperature:      -4.2,
		MaxTemperature:      25.13,
		FillLikeTemperature: 4.12,
		IconResource:        "1",
	})

	assert.Nil(t, err)
}

func TestGetCurrentWeatherEntity(t *testing.T) {
	r := initRedisCurrentRepository(t)
	TestSaveCurrentWeatherEntity(t)

	currentWeatherEntity, err := r.Get(cityName)

	assert.Nil(t, err)
	assert.NotNil(t, currentWeatherEntity)
}

func initRedisCurrentRepository(t *testing.T) *redis2.RedisCurrentWeatherRepository {
	cfg := config.ParseEnv()
	if cfg.Redis.Address == "" {
		t.Skip("skip test because redis not set")
	}

	log := logger.NewLogger(cfg.Logger)
	ctx := context.Background()
	redisClient := redis.NewClient(&ctx, cfg.Redis)
	return redis2.NewRedisCurrentWeatherRepository(log, redisClient, cfg.ExpirationDuration)
}
