package job

import (
	"context"
	"log/slog"
	"time"
)

var _ JobWithCleanUp = (*recurringJob)(nil)

func NewRecurring(name string, recurringDuration time.Duration, doWork Work) JobWithCleanUp {
	return &recurringJob{
		name:           name,
		tickerDuration: recurringDuration,
		doWork:         doWork,
		cleanUp:        noopCleanUp,
	}
}

type recurringJob struct {
	name           string
	tickerDuration time.Duration
	doWork         Work
	cleanUp        cleanUp
}

func (j *recurringJob) Name() string {
	return j.name
}

func (j *recurringJob) Run(ctx context.Context, logger *slog.Logger) {
	ticker := time.NewTicker(j.tickerDuration)

	// ticker by default only starts ticking after specified period
	// using for syntax below to force it to tick immediately
	for ; true; <-ticker.C {
		select {
		case <-ctx.Done():
			ticker.Stop()
			logger.Info("ticker stopped")

			if err := j.cleanUp.run(ctx, logger); err != nil {
				logger.Error(err.Error())
			}

			return
		default:
			if err := j.doWork(ctx, logger); err != nil {
				logger.Error(err.Error())
			}
		}
	}
}

func (j *recurringJob) CleanUpWith(work Work, timeout time.Duration) Job {
	j.cleanUp = cleanUp{
		work:           work,
		cleanUpTimeout: timeout,
	}
	return j
}
