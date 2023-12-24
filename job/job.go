package job

import (
	"context"
	"log/slog"
	"time"
)

type Work = func(ctx context.Context, logger *slog.Logger) error

type Job interface {
	Name() string
	Run(ctx context.Context, logger *slog.Logger)
}

type JobWithCleanUp interface {
	Job
	CleanUpWith(cleanUp Work, timeout time.Duration) Job
}
