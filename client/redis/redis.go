package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"weather-api/config"
)

type Client interface {
	HSet(ctx context.Context, key string, innerField string, value interface{}) error
	HGet(ctx context.Context, key string, innerField string) ([]byte, error)
	HGetAll(ctx context.Context, key string) (map[string][]byte, error)
}

type ClientImpl struct {
	redisClient *redis.Client
}

func NewClient(cfg *config.Redis) Client {
	return &ClientImpl{
		redisClient: redis.NewClient(&redis.Options{
			Addr: cfg.Address,
			DB:   0,
		}),
	}
}

func (c *ClientImpl) HSet(ctx context.Context, key string, innerField string, value interface{}) error {
	resp := c.redisClient.HSet(ctx, key, map[string]interface{}{innerField: value})
	return resp.Err()
}

func (c *ClientImpl) HGet(ctx context.Context, key string, innerField string) ([]byte, error) {
	resp := c.redisClient.HGet(ctx, key, innerField)
	if resp.Err() != nil {
		return nil, resp.Err()
	}

	return resp.Bytes()
}

func (c *ClientImpl) HGetAll(ctx context.Context, key string) (map[string][]byte, error) {
	resp := c.redisClient.HGetAll(ctx, key)
	if resp.Err() != nil {
		return nil, resp.Err()
	}

	mapStringBytes := map[string][]byte{}

	for key := range resp.Val() {
		mapStringBytes[key] = []byte(resp.Val()[key])
	}

	return mapStringBytes, resp.Err()
}
