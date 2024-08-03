package storage

import (
	"context"
	"fmt"
	"weather-api/client/redis"
	"weather-api/dto"
)

type RedisLocationStorage struct {
	client *redis.Client
}

func NewRedisLocationStorage(client *redis.Client) LocationStorage {
	return &RedisLocationStorage{
		client: client,
	}
}

const (
	keyLastTimeUseCoords = "weather-api:coords-last-time" // coords in string, like: "-10.1234,13.1948" (latitude, longitude)
	keyCoords            = "weather-api:coords"           //
	keyLocation          = "weather-api:location"         // address hash
)

func (r RedisLocationStorage) AddNewCoords(ctx context.Context, coords *dto.Coords, addressHash string) error {
	return r.client.HSet(ctx, keyCoords, r.getStrCoords(coords), addressHash)
}

func (r RedisLocationStorage) GetLastTimeGetAddressHash(ctx context.Context, coords *dto.Coords) (int64, error) {
	return getInt64WithInnerFieldsFromRedis(ctx, r.client, keyLastTimeUseCoords, r.getStrCoords(coords))
}

func (r RedisLocationStorage) UpdateLastTimeGetAddressHash(ctx context.Context, coords *dto.Coords, lastTime int64) error {
	return r.client.HSet(ctx, keyLastTimeUseCoords, r.getStrCoords(coords), lastTime)
}

func (r RedisLocationStorage) GetAddressHashByCoords(ctx context.Context, coords *dto.Coords) (string, error) {
	return getStringWithInnerFieldsFromRedis(ctx, r.client, keyCoords, r.getStrCoords(coords))
}

func (r RedisLocationStorage) GetLocation(ctx context.Context, addressHash string) (*dto.Location, error) {
	locationDto := &dto.Location{}
	err := getObjWithInnerFieldFromRedis[dto.Location](ctx, r.client, keyLocation, addressHash, locationDto)
	return locationDto, err
}

func (r RedisLocationStorage) SaveLocation(ctx context.Context, location dto.Location, addressHash string) error {
	return setObjWithInnerFieldToRedis(ctx, r.client, keyLocation, addressHash, location)
}

func (r RedisLocationStorage) getStrCoords(coords *dto.Coords) string {
	return fmt.Sprintf("%f,%f", coords.Latitude, coords.Longitude)
}
