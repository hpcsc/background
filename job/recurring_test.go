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

func TestRecurringJob(t *testing.T) {
	t.Run("return without blocking when main context is canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		j := job.WithRecurringFunc("test", 5*time.Second, func(logger *slog.Logger) {})

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
		j := job.WithRecurringWork("test", 5*time.Second, work)

		jobCtx := test.RunJobInGoroutine(j, ctx)

		// simulate shutdown signal
		cancel()

		// wait for job goroutine to finish
		<-jobCtx.Done()

		require.True(t, work.CleanUpCalled)
	})

	t.Run("trigger work periodically", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		work := test.NewWork()
		j := job.WithRecurringWork("test", 2*time.Millisecond, work)

		jobCtx := test.RunJobInGoroutine(j, ctx)

		// simulate shutdown after 5ms
		time.AfterFunc(5*time.Millisecond, func() { cancel() })

		// wait for job goroutine to finish
		<-jobCtx.Done()

		// expect ticker to run 3 times (at 0ms, 2ms, 4ms)
		require.Equal(t, 3, work.RunCount)
	})
}
