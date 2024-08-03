package scheduler

import (
	"context"
	"time"
)

type Scheduler interface {
	Run(ctx context.Context)
}

func isRowNeedRemove(now int64, lastTimeUpdated int64, duration time.Duration) bool {
	return now >= lastTimeUpdated+duration.Milliseconds()
}
