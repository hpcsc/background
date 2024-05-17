//go:build unit

package job_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/hpcsc/background/job"
	"github.com/hpcsc/background/test"
	"github.com/stretchr/testify/require"
)

func TestBlockingJob(t *testing.T) {
	t.Run("return without blocking when main context is canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		j := job.WithBlockingFunc("test", func(logger *slog.Logger) {})

		jobCtx := test.RunJobInGoroutine(j, ctx)

		// simulate shutdown signal
		cancel()

		// wait for job goroutine to finish
		<-jobCtx.Done()

		// expect job goroutine to finish before timeout
		require.NotEqual(t, context.DeadlineExceeded, jobCtx.Err())
	})

	t.Run("run cleanup logic when work implements CleanUp", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		work := test.NewWork()
		j := job.WithBlockingWork("test", work)

		jobCtx := test.RunJobInGoroutine(j, ctx)

		// simulate shutdown after 3ms
		time.AfterFunc(3*time.Millisecond, func() { cancel() })

		// wait for job goroutine to finish
		<-jobCtx.Done()

		// work is called
		require.Equal(t, 1, work.RunCount)

		// cleanup is called
		require.True(t, work.CleanUpCalled)
	})
}
