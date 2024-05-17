package job

import (
	"context"
	"log/slog"
)

func WithBlockingWork(name string, work Work) Job {
	return &blockingJob{
		name: name,
		work: work,
	}
}

func WithBlockingFunc(name string, workFunc func(logger *slog.Logger)) Job {
	return &blockingJob{
		name: name,
		work: &functionalWork{workFunc: workFunc},
	}
}

var _ Job = (*blockingJob)(nil)

type blockingJob struct {
	name string
	work Work
	workCleanUpRunner
}

func (j *blockingJob) Name() string {
	return j.name
}

func (j *blockingJob) Run(ctx context.Context, logger *slog.Logger) {
	go func() {
		j.work.Run(logger)
	}()

	<-ctx.Done()

	j.workCleanUpRunner.run(j.work, logger)
}
