package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"weather-api/config"
)

type Client struct {
	redisClient *redis.Client
}

func NewClient(cfg *config.Redis) *Client {
	return &Client{
		redisClient: redis.NewClient(&redis.Options{
			Addr: cfg.Address,
			DB:   0,
		}),
	}
}

func (c *Client) Set(ctx context.Context, key string, value interface{}) error {
	resp := c.redisClient.Set(ctx, key, value, 0)
	return resp.Err()
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	resp := c.redisClient.Get(ctx, key)

	if resp.Err() != nil {
		return nil, resp.Err()
	}

	return resp.Bytes()
}
