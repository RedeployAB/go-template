package server

import (
	"log/slog"
	"os"
)

// logger is the interface that wraps around methods Info and Error.
type logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

// NewLogger creates a new slog with a JSON handler.
func NewLogger() logger {
	return slog.New(slog.NewJSONHandler(os.Stderr, nil))
}
