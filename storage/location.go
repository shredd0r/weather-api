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
	keyLastTimeUseCoords = "weather-api:coords-last-time:%s" // coords in string, like: "-10.1234,13.1948" (latitude, longitude)
	keyCoordsFormat      = "weather-api:coords:%s"           //
	keyLocationFormat    = "weather-api:location:%s"         // address hash
)

func (r RedisLocationStorage) AddNewCoords(ctx context.Context, coords *dto.Coords, addressHash string) error {
	return r.client.Set(ctx, r.getKeyWithCoords(keyCoordsFormat, coords), addressHash)
}

func (r RedisLocationStorage) GetLastTimeGetAddressHash(ctx context.Context, coords *dto.Coords) (int64, error) {
	return getInt64FromRedis(ctx, r.client, r.getKeyWithCoords(keyLastTimeUseCoords, coords))
}

func (r RedisLocationStorage) UpdateLastTimeGetAddressHash(ctx context.Context, coords *dto.Coords, lastTime int64) error {
	return r.client.Set(ctx, r.getKeyWithCoords(keyLastTimeUseCoords, coords), lastTime)
}

func (r RedisLocationStorage) GetAddressHashByCoords(ctx context.Context, coords *dto.Coords) (string, error) {
	return getStringFromRedis(ctx, r.client, r.getKeyWithCoords(keyCoordsFormat, coords))
}

func (r RedisLocationStorage) GetLocation(ctx context.Context, addressHash string) (*dto.Location, error) {
	locationDto := &dto.Location{}
	err := getObjFromRedis[dto.Location](ctx, r.client, fmt.Sprintf(keyLocationFormat, addressHash), locationDto)
	return locationDto, err
}

func (r RedisLocationStorage) SaveLocation(ctx context.Context, location dto.Location, addressHash string) error {
	return setObjToRedis(ctx, r.client, fmt.Sprintf(keyLocationFormat, addressHash), location)
}

func (r RedisLocationStorage) getKeyWithCoords(format string, coords *dto.Coords) string {
	strCoords := fmt.Sprintf("%f,%f", coords.Latitude, coords.Longitude)
	return fmt.Sprintf(format, strCoords)
}
