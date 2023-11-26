package job

import (
	"context"
	"log/slog"
	"time"
)

var _ Interface = (*tickerJob)(nil)

func NewTicker(name string, tickerDuration time.Duration, doWork Work) Interface {
	return &tickerJob{
		name:           name,
		tickerDuration: tickerDuration,
		doWork:         doWork,
	}
}

type tickerJob struct {
	name           string
	tickerDuration time.Duration
	doWork         Work
}

func (j *tickerJob) Name() string {
	return j.name
}

func (j *tickerJob) Run(ctx context.Context, logger *slog.Logger) {
	ticker := time.NewTicker(j.tickerDuration)

	// ticker by default only starts ticking after specified period
	// using for syntax below to force it to tick immediately
	for ; true; <-ticker.C {
		select {
		case <-ctx.Done():
			ticker.Stop()
			logger.Info("ticker stopped")

			return
		default:
			if err := j.doWork(ctx, logger); err != nil {
				logger.Error(err.Error())
			}
		}
	}
}
