package job

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

var _ JobWithCleanUp = (*blockingJob)(nil)

func NewBlocking(name string, runAndBlock Work) JobWithCleanUp {
	return &blockingJob{
		name:        name,
		runAndBlock: runAndBlock,
		cleanUp:     noopCleanUp,
	}
}

type blockingJob struct {
	name        string
	runAndBlock Work
	cleanUp     cleanUp
}

func (j *blockingJob) Name() string {
	return j.name
}

func (j *blockingJob) Run(ctx context.Context, logger *slog.Logger) {
	// `ctx` is used to communicate whether the application is canceled due to signals
	// create a new context `runCtx` for clean up goroutine to communicate with this goroutine
	runCtx, cancelRunCtx := context.WithCancel(context.Background())

	go func() {
		// wait for shutdown signal
		<-ctx.Done()

		if err := j.cleanUp.run(ctx, logger); err != nil {
			logger.Error(err.Error())
			return
		}

		// clean up done, notify run context that it's done
		cancelRunCtx()
	}()

	if err := j.runAndBlock(ctx, logger); err != nil {
		logger.Error(fmt.Sprintf("failed to run: %v", err))
	}

	// wait for run context to be stopped
	<-runCtx.Done()
}

func (j *blockingJob) CleanUpWith(work Work, timeout time.Duration) Job {
	j.cleanUp = cleanUp{
		work:           work,
		cleanUpTimeout: timeout,
	}
	return j
}
