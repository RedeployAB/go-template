package service

import (
	"log/slog"
	"os"
)

// logger is the interface that wraps around methods Info and Error.
type logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

// NewDefaultLogger creates a new default logger.
func NewDefaultLogger() logger {
	return slog.New(slog.NewJSONHandler(os.Stderr, nil))
}
