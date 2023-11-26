package job

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

var _ Interface = (*blockingJob)(nil)

func NewBlocking(name string, runAndBlock Work, cleanup Work) Interface {
	return &blockingJob{
		name:           name,
		runAndBlock:    runAndBlock,
		cleanup:        cleanup,
		cleanupTimeout: 30 * time.Second,
	}
}

type blockingJob struct {
	name           string
	runAndBlock    Work
	cleanup        Work
	cleanupTimeout time.Duration
}

func (j *blockingJob) WithCleanupTimeout(timeout time.Duration) *blockingJob {
	j.cleanupTimeout = timeout
	return j
}

func (j *blockingJob) Name() string {
	return j.name
}

func (j *blockingJob) Run(ctx context.Context, logger *slog.Logger) {
	// `ctx` is used to communicate whether the application is canceled due to signals
	mainCtx, cancelMainCtx := context.WithCancel(context.Background())

	go func() {
		// wait for shutdown signal
		<-ctx.Done()

		// allow fixed timeout for clean up
		cleanupCtx, cancelCleanupTimeout := context.WithTimeout(mainCtx, j.cleanupTimeout)

		// cancel clean up timeout if clean up is done before timeout is reached
		defer cancelCleanupTimeout()

		go func() {
			// wait for cleanup
			<-cleanupCtx.Done()

			if errors.Is(cleanupCtx.Err(), context.DeadlineExceeded) {
				// clean up doesn't finish within timeout period, panic so that the goroutine is forced to exit
				panic(fmt.Errorf("cleanup timed out after %v", j.cleanupTimeout))
			}
		}()

		// trigger clean up and block
		logger.Info("cleaning up")
		if err := j.cleanup(cleanupCtx, logger); err != nil {
			// something is wrong with the clean up, panic to force goroutine to shut down
			panic(fmt.Errorf("failed to run clean up: %v", err))
		}

		// clean up done, notify main context that it's done
		cancelMainCtx()

		// the defer above will cancel the clean up timeout
	}()

	if err := j.runAndBlock(mainCtx, logger); err != nil {
		logger.Error(fmt.Sprintf("failed to run: %v", err))
	}

	// wait for server context to be stopped
	<-mainCtx.Done()
}
