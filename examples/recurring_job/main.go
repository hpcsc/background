package main

import (
	"context"
	"github.com/hpcsc/background"
	"github.com/hpcsc/background/job"
	"log/slog"
	"os"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	runner := background.NewRunner(context.Background(), logger)

	runner.Run(job.NewRecurring("job-1", 3*time.Second, func(_ context.Context, l *slog.Logger) error {
		l.Info("processing")
		return nil
	}))

	jobWithCleanUpLogic := job.NewRecurring("job-2", 5*time.Second, func(_ context.Context, l *slog.Logger) error {
		l.Info("processing")
		return nil
	}).
		CleanUpWith(func(_ context.Context, l *slog.Logger) error {
			l.Info("custom clean up logic")
			return nil
		}, 30*time.Second)

	runner.Run(jobWithCleanUpLogic)

	runner.Wait()

	logger.Info("exit")
}
