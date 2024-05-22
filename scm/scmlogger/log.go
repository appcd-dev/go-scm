package scmlogger

import (
	"context"
	"log/slog"
)

type LoggerFunc func(ctx context.Context) *slog.Logger

var loggerFunc LoggerFunc

func SetLoggerFunc(f LoggerFunc) {
	loggerFunc = f
}

func GetLogger(ctx context.Context) *slog.Logger {
	if loggerFunc == nil {
		return slog.Default()
	}
	return loggerFunc(ctx).With("namespace", "scm")
}

type noopLogger struct{}

func (l *noopLogger) Log(_ string, _ ...interface{}) {}
