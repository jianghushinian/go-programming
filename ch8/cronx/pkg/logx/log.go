package logx

import (
	"context"
	"log/slog"
	"strings"
	"sync"
)

var (
	// defaultLogger is the global shared logger instance.
	defaultLogger *slog.Logger
	// initOnce ensures Init() executes only once.
	initOnce sync.Once
)

// Init initializes the global logger using the provided Options.
// This function is safe to call multiple times, but initialization
// will only happen once.
func Init(opt *Options) {
	initOnce.Do(func() {
		defaultLogger = newLogger(opt)
	})
}

// Sync is kept for API compatibility.
// slog does not provide a flush mechanism, so this is a no-op.
func Sync() {
	// no-op
}

// newLogger builds a slog.Logger using the provided Options.
// It configures log level, format, output paths, and wraps
// the handler with trace/job ID injection support.
func newLogger(opt *Options) *slog.Logger {
	var level slog.Level
	switch strings.ToLower(opt.Level) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handler := buildHandler(opt, level)
	return slog.New(handler)
}

// Debug logs a debug-level message with context.
func Debug(ctx context.Context, msg string, args ...any) {
	defaultLogger.DebugContext(ctx, msg, args...)
}

// Info logs an info-level message with context.
func Info(ctx context.Context, msg string, args ...any) {
	defaultLogger.InfoContext(ctx, msg, args...)
}

// Warn logs a warning-level message with context.
func Warn(ctx context.Context, msg string, args ...any) {
	defaultLogger.WarnContext(ctx, msg, args...)
}

// Error logs an error-level message with context.
func Error(ctx context.Context, msg string, args ...any) {
	defaultLogger.ErrorContext(ctx, msg, args...)
}

// L returns the global default logger instance.
func L() *slog.Logger {
	return defaultLogger
}
