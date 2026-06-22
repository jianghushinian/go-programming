package cronx

import (
	"context"

	"cronx/pkg/logx"
)

// cronLogger implement the cron.Logger interface.
type cronLogger struct{}

// newCronLogger returns a cron logger.
func newCronLogger() *cronLogger {
	return &cronLogger{}
}

// Info logs routine messages about cron's operation.
func (l *cronLogger) Info(msg string, keysAndValues ...any) {
	logx.Info(context.Background(), msg, keysAndValues...)
}

// Error logs an error condition.
func (l *cronLogger) Error(err error, msg string, keysAndValues ...any) {
	logx.Error(context.Background(), msg, append(keysAndValues, "err", err)...)
}
