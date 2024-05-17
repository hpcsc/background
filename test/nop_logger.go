//go:build unit

package test

import (
	"context"
	"log/slog"
)

func NewNopLogger() *slog.Logger {
	return slog.New(nopLogHandler{})
}

type nopLogHandler struct {
}

func (h nopLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return false
}

func (h nopLogHandler) Handle(ctx context.Context, record slog.Record) error {
	return nil
}

func (h nopLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h nopLogHandler) WithGroup(name string) slog.Handler {
	return h
}
