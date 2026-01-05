// Package slogdiscard provides a discard logger handler that ignores all log entries.
package slogdiscard

import (
	"context"
	"log/slog"
)

// NewDiscardLogger creates a new logger that discards all log entries.
func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

// DiscardHandler is a slog handler that discards all log records.
type DiscardHandler struct{}

// NewDiscardHandler creates a new discard handler.
func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

// Handle discards the log record.
func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

// WithAttrs returns the same handler since attributes are not stored.
func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup returns the same handler since groups are not stored.
func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	return h
}

// Enabled always returns false since log entries are ignored.
func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
