//go:build unit

package test

import (
	"context"
	"log/slog"

	"github.com/hpcsc/background/job"
)

var _ job.Job = (*Job)(nil)

func NewJob() *Job {
	return &Job{}
}

type Job struct {
}

func (j *Job) Name() string {
	return "test"
}

func (j *Job) Run(ctx context.Context, logger *slog.Logger) {
}
