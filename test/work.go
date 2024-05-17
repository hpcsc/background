//go:build unit

package test

import (
	"log/slog"
	"time"

	"github.com/hpcsc/background/job"
)

func NewWork() *WorkWithCleanUp {
	return &WorkWithCleanUp{}
}

var _ job.WorkWithCleanUp = (*WorkWithCleanUp)(nil)

type WorkWithCleanUp struct {
	RunCount      int
	CleanUpCalled bool
}

func (w *WorkWithCleanUp) Run(logger *slog.Logger) {
	w.RunCount = w.RunCount + 1
}

func (w *WorkWithCleanUp) CleanUp() {
	w.CleanUpCalled = true
}

func (w *WorkWithCleanUp) CleanUpTimeOut() time.Duration {
	return 2 * time.Second
}
