package storage

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	goredis "github.com/redis/go-redis/v9"
	"weather-api/client/redis"
)

func setObjToRedis[T any](ctx context.Context, client *redis.Client, key string, value T) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return client.Set(ctx, key, bytes)
}

func getObjFromRedis[T any](ctx context.Context, client *redis.Client, key string, valueResp *T) error {

	bytes, err := getBytesFromRedis(ctx, client, key)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, valueResp)

	if err != nil {
		return err
	}

	return nil
}

func getStringFromRedis(ctx context.Context, client *redis.Client, key string) (string, error) {
	bytes, err := getBytesFromRedis(ctx, client, key)
	if err != nil {
		return "", err
	}

	return string(bytes), err
}

func getInt64FromRedis(ctx context.Context, client *redis.Client, key string) (int64, error) {
	bytes, err := getBytesFromRedis(ctx, client, key)
	if err != nil {
		return -1, err
	}

	return int64(binary.BigEndian.Uint64(bytes)), err
}

func getBytesFromRedis(ctx context.Context, client *redis.Client, key string) ([]byte, error) {
	bytes, err := client.Get(ctx, key)

	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return bytes, nil
}
