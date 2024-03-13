package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
	"weather-api/config"
)

type Client struct {
	ctx         *context.Context
	redisClient *redis.Client
}

func NewClient(ctx *context.Context, cfg config.Redis) *Client {
	return &Client{
		ctx: ctx,
		redisClient: redis.NewClient(
			&redis.Options{
				Addr: cfg.Address,
			}),
	}
}

func (c *Client) Set(key string, value string, expiration time.Duration) error {
	return c.redisClient.Set(*c.ctx, key, value, expiration).Err()
}

func (c *Client) Get(key string) (string, error) {
	return c.redisClient.Get(*c.ctx, key).Result()
}
