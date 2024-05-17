package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/hpcsc/background"
	"github.com/hpcsc/background/job"
)

var _ job.WorkWithCleanUp = (*recurringWorkWithCleanUp)(nil)

type recurringWorkWithCleanUp struct {
	logger *slog.Logger
}

func (j *recurringWorkWithCleanUp) Run(logger *slog.Logger) {
	logger.Info("processing")
}

func (j *recurringWorkWithCleanUp) CleanUp() {
	j.logger.Info("cleaning up")
}

func (j *recurringWorkWithCleanUp) CleanUpTimeOut() time.Duration {
	return 5 * time.Second
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	runner := background.NewRunner(context.Background(), logger)

	runner.Run(
		job.WithRecurringFunc(
			"job-1",
			3*time.Second,
			func(l *slog.Logger) {
				l.Info("processing")
			},
		),
	)

	runner.Run(
		job.WithRecurringWork(
			"job-2",
			5*time.Second,
			&recurringWorkWithCleanUp{logger: logger},
		),
	)

	runner.Wait()

	logger.Info("exit")
}
