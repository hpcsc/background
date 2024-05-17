package job

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type workCleanUpRunner struct{}

func (r workCleanUpRunner) run(work Work, logger *slog.Logger) {
	if cleanUp, ok := work.(CleanUp); ok {
		logger.Info("waiting for cleaning up")
		r.runCleanUp(cleanUp, logger)
		logger.Info("all clean up works done")
	} else {
		logger.Info("no clean up needed")
	}
}

func (r workCleanUpRunner) runCleanUp(cleanUp CleanUp, logger *slog.Logger) {
	// allow fixed timeout for clean up
	timeOut := cleanUp.CleanUpTimeOut()
	cleanupCtx, cancelCleanupTimeout := context.WithTimeout(context.Background(), timeOut)

	// cancel clean up timeout if clean up is done before timeout is reached
	defer cancelCleanupTimeout()

	go func() {
		// wait for cleanup
		<-cleanupCtx.Done()

		if errors.Is(cleanupCtx.Err(), context.DeadlineExceeded) {
			// clean up doesn't finish within timeout period, panic so that the goroutine is forced to exit
			logger.Error(fmt.Sprintf("clean up timed out after %v", timeOut))
		}
	}()

	// trigger clean up
	cleanUp.CleanUp()
}
