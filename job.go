package background

import (
	"context"
	"log/slog"
)

type Work = func(ctx context.Context, logger *slog.Logger) error

type Job interface {
	Name() string
	Run(ctx context.Context, logger *slog.Logger)
}
