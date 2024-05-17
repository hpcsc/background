package job

import "log/slog"

var _ Work = (*functionalWork)(nil)

type functionalWork struct {
	workFunc func(logger *slog.Logger)
}

func (w *functionalWork) Run(logger *slog.Logger) {
	w.workFunc(logger)
}
