package storage

import (
	"context"
	"errors"
	"fmt"

	redis2 "github.com/redis/go-redis/v9"
	"github.com/shredd0r/weather-api/client/redis"
	"github.com/shredd0r/weather-api/dto"
	"github.com/shredd0r/weather-api/log"
	"github.com/shredd0r/weather-api/util"
)

const (
	keyLastTimeUseCoords = "weather-api:coords-last-time" // coords in string, like: "-10.1234,13.1948" (latitude, longitude)
	keyCoords            = "weather-api:coords"           //
	keyLocation          = "weather-api:location"         // address hash
)

type RedisLocationStorage struct {
	logger log.Logger
	client redis.Client
}

func NewRedisLocationStorage(logger log.Logger, client redis.Client) LocationStorage {
	return &RedisLocationStorage{
		logger: logger,
		client: client,
	}
}

func (r RedisLocationStorage) GetLocation(ctx context.Context, addressHash string) (*dto.Location, error) {
	r.logger.Debug("try get location from redis cache")
	locationDto := &dto.Location{}
	err := getObjWithInnerFieldFromRedis[dto.Location](ctx, r.client, keyLocation, addressHash, locationDto)
	return locationDto, err
}

func (r RedisLocationStorage) SaveLocation(ctx context.Context, location *dto.Location, addressHash string) error {
	r.logger.Debug("try save new location to redis cache")
	return setObjWithInnerFieldToRedis(ctx, r.client, keyLocation, addressHash, location)
}

func (r RedisLocationStorage) AddCoords(ctx context.Context, coords *dto.Coords, addressHash string) error {
	r.logger.Debugf("try add new coords to redis cache, coords: %s, addressHash: %s", coords, addressHash)
	return r.client.HSet(ctx, keyCoords, r.getStrCoords(coords), addressHash)
}

func (r RedisLocationStorage) RemoveCoords(ctx context.Context, coords *dto.Coords) error {
	r.logger.Debugf("try delete coords from redis cache, coords: %s", coords)
	return r.client.HDel(ctx, keyLastTimeUseCoords, r.getStrCoords(coords))
}

func (r RedisLocationStorage) UpdateLastTimeUseCoords(ctx context.Context, coords *dto.Coords, lastTime int64) error {
	r.logger.Debugf("try update last time use coords in redis cache, coords: %s, lastTime: %d", coords, lastTime)
	return r.client.HSet(ctx, keyLastTimeUseCoords, r.getStrCoords(coords), lastTime)
}

func (r RedisLocationStorage) GetAllLastTimeUseCoords(ctx context.Context) (map[dto.Coords]int64, error) {
	r.logger.Debugf("try get all last time use coords in redis cache")

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
	r.logger.Debugf("try get addressHash by coords from redis cache")
	return getStringWithInnerFieldsFromRedis(ctx, r.client, keyCoords, r.getStrCoords(coords))
}

func (r RedisLocationStorage) getStrCoords(coords *dto.Coords) string {
	return fmt.Sprintf("%f,%f", coords.Latitude, coords.Longitude)
}
