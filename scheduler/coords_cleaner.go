package scheduler

import (
	"context"
	"sync"
	"time"
	"weather-api/config"
	"weather-api/dto"
	"weather-api/log"
	"weather-api/storage"
)

type CoordsCleaner struct {
	logger  log.Logger
	wg      *sync.WaitGroup
	cfg     *config.ExpirationDuration
	storage storage.LocationStorage
}

func (c *CoordsCleaner) Run(ctx context.Context) {
	c.workFlowForRemoveCoordsFromCache(ctx)
}

func (c *CoordsCleaner) workFlowForRemoveCoordsFromCache(ctx context.Context) {
	c.wg.Add(1)

	go func() {
		defer c.wg.Done()

		mapCoordsToLastTime, err := c.storage.GetAllLastTimeUseCoords(ctx)
		if err != nil {
			c.logger.Error("error when try get all last time use coords")
		}

		now := time.Now().UnixMilli()
		for coords, lastTime := range mapCoordsToLastTime {
			if isRowNeedRemove(now, lastTime, c.cfg.Coords) {
				c.removeCoordsInGoroutine(ctx, &coords)
			}
		}

	}()

	c.wg.Wait()
}

func (c *CoordsCleaner) removeCoordsInGoroutine(ctx context.Context, coords *dto.Coords) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		err := c.storage.RemoveCoords(ctx, coords)

		if err != nil {
			c.logger.Error("error when try remove coords from cache")
		}
	}()
}
