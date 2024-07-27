package storage

import (
	"context"
	"github.com/redis/go-redis/v9"
	"weather-api/dto"
)

type RedisLocationStorage struct {
	client redis.Client
}

func (r RedisLocationStorage) AddNewCoords(ctx context.Context, coords *dto.Coords, addressHash string) error {
	//TODO implement me
	return nil
}

func (r RedisLocationStorage) LastTimeGetAddressHash(ctx context.Context, coords *dto.Coords) (uint64, error) {
	//TODO implement me
	return 0, nil
}

func (r RedisLocationStorage) UpdateLastTimeGetAddressHash(ctx context.Context, coords *dto.Coords) error {
	//TODO implement me
	return nil
}

func (r RedisLocationStorage) GetAddressHashByCoords(ctx context.Context, coords *dto.Coords) (string, error) {
	//TODO implement me
	return "", ErrNotFound
}

func (r RedisLocationStorage) GetLocation(ctx context.Context, addressHash string) (*dto.Location, error) {
	//TODO implement me
	return nil, ErrNotFound
}

func (r RedisLocationStorage) SaveLocation(ctx context.Context, location *dto.Location) error {
	//TODO implement me
	return nil
}
