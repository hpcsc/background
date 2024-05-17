package job

import (
	"log/slog"
	"time"
)

type Work interface {
	Run(logger *slog.Logger)
}

type CleanUp interface {
	CleanUp()
	CleanUpTimeOut() time.Duration
}

type WorkWithCleanUp interface {
	Work
	CleanUp
}
