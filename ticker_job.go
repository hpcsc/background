package background

import (
	"context"
	"log/slog"
	"time"
)

var _ Job = (*tickerJob)(nil)

func NewTickerJob(name string, tickerDuration time.Duration, runMain func(ctx context.Context, logger *slog.Logger)) Job {
	return &tickerJob{
		name:           name,
		tickerDuration: tickerDuration,
		runMain:        runMain,
	}
}

type tickerJob struct {
	name           string
	tickerDuration time.Duration
	runMain        func(ctx context.Context, logger *slog.Logger)
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
			j.runMain(ctx, logger)
		}
	}
}
