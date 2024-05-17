package job

import (
	"context"
	"log/slog"
)

type Job interface {
	Name() string
	Run(ctx context.Context, logger *slog.Logger)
}
