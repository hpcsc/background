package job

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

const defaultCleanUpTimeout = 30 * time.Second

var noopCleanUp = cleanUp{
	work: func(_ context.Context, _ *slog.Logger) error {
		return nil
	},
	cleanUpTimeout: defaultCleanUpTimeout,
}

type cleanUp struct {
	work           Work
	cleanUpTimeout time.Duration
}

func (c *cleanUp) run(ctx context.Context, logger *slog.Logger) error {
	// allow fixed timeout for clean up
	cleanupCtx, cancelCleanupTimeout := context.WithTimeout(ctx, c.cleanUpTimeout)

	// cancel clean up timeout if clean up is done before timeout is reached
	defer cancelCleanupTimeout()

	go func() {
		// wait for cleanup
		<-cleanupCtx.Done()

		if errors.Is(cleanupCtx.Err(), context.DeadlineExceeded) {
			// clean up doesn't finish within timeout period, panic so that the goroutine is forced to exit
			logger.Error(fmt.Sprintf("cleanup timed out after %v", c.cleanUpTimeout))
		}
	}()

	// trigger clean up
	if err := c.work(cleanupCtx, logger); err != nil {
		// something is wrong with the clean up
		return fmt.Errorf("failed to run clean up: %v", err)
	}

	return nil
}
