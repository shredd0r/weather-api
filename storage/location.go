package storage

import (
	"context"
	"errors"
	"fmt"
	redis2 "github.com/redis/go-redis/v9"
	"weather-api/client/redis"
	"weather-api/dto"
	"weather-api/util"
)

const (
	keyLastTimeUseCoords = "weather-api:coords-last-time" // coords in string, like: "-10.1234,13.1948" (latitude, longitude)
	keyCoords            = "weather-api:coords"           //
	keyLocation          = "weather-api:location"         // address hash
)

type RedisLocationStorage struct {
	client redis.Client
}

func NewRedisLocationStorage(client redis.Client) LocationStorage {
	return &RedisLocationStorage{
		client: client,
	}
}

func (r RedisLocationStorage) GetLocation(ctx context.Context, addressHash string) (*dto.Location, error) {
	locationDto := &dto.Location{}
	err := getObjWithInnerFieldFromRedis[dto.Location](ctx, r.client, keyLocation, addressHash, locationDto)
	return locationDto, err
}

func (r RedisLocationStorage) SaveLocation(ctx context.Context, location dto.Location, addressHash string) error {
	return setObjWithInnerFieldToRedis(ctx, r.client, keyLocation, addressHash, location)
}

func (r RedisLocationStorage) AddCoords(ctx context.Context, coords *dto.Coords, addressHash string) error {
	return r.client.HSet(ctx, keyCoords, r.getStrCoords(coords), addressHash)
}

func (r RedisLocationStorage) RemoveCoords(ctx context.Context, coords *dto.Coords) error {
	return nil
}

func (r RedisLocationStorage) UpdateLastTimeUseCoords(ctx context.Context, coords *dto.Coords, lastTime int64) error {
	return r.client.HSet(ctx, keyLastTimeUseCoords, r.getStrCoords(coords), lastTime)
}

func (r RedisLocationStorage) GetAllLastTimeUseCoords(ctx context.Context) (map[dto.Coords]int64, error) {
	mapStrCoordsToLastTimeUse, err := r.client.HGetAll(ctx, keyLastTimeUseCoords)
	if err != nil {
		if errors.Is(err, redis2.Nil) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	mapCoordsToLastTimeUse := map[dto.Coords]int64{}
	for strCoords, lastTime := range mapStrCoordsToLastTimeUse {
		mapCoordsToLastTimeUse[dto.NewCoords(strCoords)] = util.BytesToInt64(lastTime)
	}

	return mapCoordsToLastTimeUse, nil
}

func (r RedisLocationStorage) GetAddressHashByCoords(ctx context.Context, coords *dto.Coords) (string, error) {
	return getStringWithInnerFieldsFromRedis(ctx, r.client, keyCoords, r.getStrCoords(coords))
}

func (r RedisLocationStorage) getStrCoords(coords *dto.Coords) string {
	return fmt.Sprintf("%f,%f", coords.Latitude, coords.Longitude)
}
