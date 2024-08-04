package task

import (
	"context"
	"github.com/shredd0r/weather-api/config"
	"github.com/shredd0r/weather-api/dto"
	"github.com/shredd0r/weather-api/log"
	"github.com/shredd0r/weather-api/storage"
	"sync"
	"time"
)

type CoordsCleaner struct {
	logger  log.Logger
	wg      *sync.WaitGroup
	cfg     *config.ExpirationDuration
	storage storage.LocationStorage
}

func NewCoordsCleaner(logger log.Logger, cfg *config.ExpirationDuration, storage storage.LocationStorage) *CoordsCleaner {
	return &CoordsCleaner{
		logger:  logger,
		wg:      &sync.WaitGroup{},
		cfg:     cfg,
		storage: storage,
	}
}

func (c *CoordsCleaner) Run(ctx context.Context) {
	c.workFlowForRemoveCoordsFromCache(ctx)
}

func (c *CoordsCleaner) workFlowForRemoveCoordsFromCache(ctx context.Context) {
	c.logger.Debug("start remove coords from cache")
	mapCoordsToLastTime, err := c.storage.GetAllLastTimeUseCoords(ctx)
	if err != nil {
		c.logger.Error("error when try get all last time use coords")
	}

	c.logger.Debugf("fined %d coords in cache", len(mapCoordsToLastTime))

	now := time.Now().UnixMilli()
	for coords, lastTime := range mapCoordsToLastTime {
		if isRowNeedRemove(now, lastTime, c.cfg.Coords) {
			c.logger.Debugf("coords %s need remove", coords)
			c.removeCoordsInGoroutine(ctx, &coords)
		}
	}

}

func (c *CoordsCleaner) removeCoordsInGoroutine(ctx context.Context, coords *dto.Coords) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		err := c.storage.RemoveCoords(ctx, coords)

		if err != nil {
			c.logger.Error("error when try remove coords from cache")
		}

		c.logger.Debugf("coords %s removed", coords)
	}()
}
