package job

import (
	"context"
	"log/slog"
	"time"
)

func WithRecurringWork(name string, interval time.Duration, work Work) Job {
	return &recurringJob{
		name:     name,
		interval: interval,
		work:     work,
	}
}

func WithRecurringFunc(name string, interval time.Duration, workFunc func(logger *slog.Logger)) Job {
	return &recurringJob{
		name:     name,
		interval: interval,
		work:     &functionalWork{workFunc: workFunc},
	}
}

var _ Job = (*recurringJob)(nil)

type recurringJob struct {
	name     string
	interval time.Duration
	work     Work
	workCleanUpRunner
}

func (j *recurringJob) Name() string {
	return j.name
}

func (j *recurringJob) Run(ctx context.Context, logger *slog.Logger) {
	ticker := time.NewTicker(j.interval)

	go func() {
		// ticker by default only starts ticking after specified period
		// using for syntax below to force it to tick immediately
		for ; true; <-ticker.C {
			j.work.Run(logger)
		}
	}()

	<-ctx.Done()

	ticker.Stop()

	j.workCleanUpRunner.run(j.work, logger)
}
