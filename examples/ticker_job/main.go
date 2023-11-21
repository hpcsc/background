package main

import (
	"context"
	"github.com/hpcsc/background"
	"log/slog"
	"os"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	runner := background.NewRunner(context.Background(), logger)

	runner.Run(background.NewTickerJob("job-1", 3*time.Second, func(_ context.Context, l *slog.Logger) error {
		l.Info("processing")
		return nil
	}))

	runner.Run(background.NewTickerJob("job-2", 5*time.Second, func(_ context.Context, l *slog.Logger) error {
		l.Info("processing")
		return nil
	}))

	runner.Wait()

	logger.Info("exit")
}
