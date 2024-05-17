package background

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hpcsc/background/job"
)

func NewRunner(parentCtx context.Context, logger *slog.Logger) *Runner {
	// let caller pass in a context so that the caller can cancel this runner any time it wants
	// wrap parent context with a cancel context so that the runner is always able to handle shutdown signal
	// all jobs run by this runner use the same wrapped cancel context above
	return &Runner{
		wg:     &sync.WaitGroup{},
		logger: logger,
		ctx:    runShutdownSignalHandler(parentCtx, logger),
	}
}

type Runner struct {
	wg     *sync.WaitGroup
	logger *slog.Logger
	ctx    context.Context
}

func (r *Runner) Run(job job.Job) {
	r.wg.Add(1)

	l := r.logger.With("job", job.Name())

	go func() {
		defer func() {
			if err := recover(); err != nil {
				l.Error(fmt.Sprintf("received unhandled error: %v", err))
			}

			r.wg.Done()
			l.Info("job stopped")
		}()

		l.Info("job started")

		job.Run(r.ctx, l)
	}()
}

func (r *Runner) Wait() {
	r.wg.Wait()
}

func runShutdownSignalHandler(parentCtx context.Context, logger *slog.Logger) context.Context {
	ctx, cancel := context.WithCancel(parentCtx)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		s := <-sig
		logger.Info(fmt.Sprintf("received %v signal", s))

		// signal to all other goroutines to start doing their cleanup
		cancel()
	}()

	return ctx
}
