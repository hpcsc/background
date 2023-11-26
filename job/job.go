package job

import (
	"context"
	"log/slog"
)

type Work = func(ctx context.Context, logger *slog.Logger) error

type Interface interface {
	Name() string
	Run(ctx context.Context, logger *slog.Logger)
}
