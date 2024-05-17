//go:build unit

package test

import (
	"context"
	"time"

	"github.com/hpcsc/background/job"
)

func RunJobInGoroutine(j job.Job, mainCtx context.Context) context.Context {
	jobCtx, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Millisecond)

	go func() {
		j.Run(mainCtx, NewNopLogger())

		// job completed, cancel timeout
		cancelTimeout()
	}()

	return jobCtx
}
