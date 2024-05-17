//go:build unit

package background_test

import (
	"context"
	"testing"
	"time"

	"github.com/hpcsc/background"
	"github.com/hpcsc/background/test"
	"github.com/stretchr/testify/require"
)

func TestRunner(t *testing.T) {
	t.Run("return without blocking when parent context is canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		r := background.NewRunner(ctx, test.NewNopLogger())
		r.Run(test.NewJob())
		r.Run(test.NewJob())

		runnerCtx := waitRunnerInGoroutine(r)

		// cancel main context
		// this can only test that runner lets all jobs finish before returning
		// it cannot test the signal handling logic
		cancel()

		// wait for runner goroutine to finish
		<-runnerCtx.Done()

		// expect wait runner goroutine to finish before timeout
		require.NotEqual(t, context.DeadlineExceeded, runnerCtx.Err())
	})
}

func waitRunnerInGoroutine(runner *background.Runner) context.Context {
	runnerCtx, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Millisecond)

	go func() {
		runner.Wait()

		// runner finished, cancel timeout
		cancelTimeout()
	}()

	return runnerCtx
}
