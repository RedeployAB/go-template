package server

import (
	"log"
	"strings"
)

// defaultLogger is a default logger.
type defaultLogger struct {
	out func(v ...any)
}

// NewDefaultLogger creates a new default logger.
func NewDefaultLogger() defaultLogger {
	return defaultLogger{
		out: log.Println,
	}
}

// Info logs an info message.
func (l defaultLogger) Info(msg string, keysAndValues ...any) {
	var b strings.Builder
	b.WriteString("message=")
	b.WriteString(msg)
	for i, line := range keysAndValues {
		if i%2 == 0 {
			b.WriteString("; ")
			b.WriteString(line.(string))
			b.WriteString("=")
		} else {
			b.WriteString(line.(string))
		}
	}
	l.out(b.String())
}

// Error logs an error message.
func (l defaultLogger) Error(err error, msg string, keysAndValues ...any) {
	var b strings.Builder
	b.WriteString("message=")
	b.WriteString(msg)
	b.WriteString(" error=")
	b.WriteString(err.Error())
	for i, line := range keysAndValues {
		if i%2 == 0 {
			b.WriteString("; ")
			b.WriteString(line.(string))
			b.WriteString("=")
		} else {
			b.WriteString(line.(string))
		}
	}
	l.out(b.String())
}
