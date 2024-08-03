package storage

import (
	"context"
	"encoding/json"
	"errors"
	goredis "github.com/redis/go-redis/v9"
	"weather-api/client/redis"
	"weather-api/util"
)

func setObjWithInnerFieldToRedis[T any](ctx context.Context, client redis.Client, key string, innerField string, value T) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return client.HSet(ctx, key, innerField, bytes)
}

func getObjWithInnerFieldFromRedis[T any](ctx context.Context, client redis.Client, key string, innerField string, valueResp *T) error {
	bytes, err := getBytesWithInnerFieldFromRedis(ctx, client, key, innerField)

	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, valueResp)

	if err != nil {
		return err
	}

	return nil
}

func getStringWithInnerFieldsFromRedis(ctx context.Context, client redis.Client, key string, innerField string) (string, error) {
	bytes, err := getBytesWithInnerFieldFromRedis(ctx, client, key, innerField)
	if err != nil {
		return "", err
	}

	return string(bytes), err
}

func getInt64WithInnerFieldsFromRedis(ctx context.Context, client redis.Client, key string, innerField string) (int64, error) {
	bytes, err := getBytesWithInnerFieldFromRedis(ctx, client, key, innerField)
	if err != nil {
		return -1, err
	}

	return util.BytesToInt64(bytes), err
}

func getBytesWithInnerFieldFromRedis(ctx context.Context, client redis.Client, key string, innerField string) ([]byte, error) {
	bytes, err := client.HGet(ctx, key, innerField)

	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return bytes, nil
}
