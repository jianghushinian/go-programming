package logx

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
)

// traceHandler is a wrapper around slog.Handler.
// It injects traceID and jobID from context into each log record.
type traceHandler struct {
	handler slog.Handler
}

func (h *traceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *traceHandler) Handle(ctx context.Context, r slog.Record) error {
	// Inject traceID from context into the log record.
	if traceID, ok := TraceIDFromContext(ctx); ok {
		r.Add("traceID", slog.StringValue(traceID))
	}
	// Inject jobID from context into the log record.
	if jobID, ok := JobIDFromContext(ctx); ok {
		r.Add("jobID", slog.Int64Value(jobID))
	}
	return h.handler.Handle(ctx, r)
}

func (h *traceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &traceHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *traceHandler) WithGroup(name string) slog.Handler {
	return &traceHandler{handler: h.handler.WithGroup(name)}
}

// buildHandler creates a slog.Handler based on Options
// and wraps it with traceHandler to support traceID/jobID injection.
func buildHandler(opt *Options, lvl slog.Level) slog.Handler {
	writer := buildWriter(opt.OutputPaths)

	var handler slog.Handler
	switch strings.ToLower(opt.Format) {
	case "json":
		handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level: lvl,
		})
	default:
		handler = slog.NewTextHandler(writer, &slog.HandlerOptions{
			Level: lvl,
		})
	}

	// Wrap with custom handler to add trace/job IDs from context.
	return &traceHandler{handler: handler}
}

// buildWriter creates an io.Writer based on the configured output paths.
// Supports multiple outputs via io.MultiWriter.
func buildWriter(paths []string) io.Writer {
	if len(paths) == 0 {
		return os.Stdout
	}

	var writers []io.Writer
	for _, p := range paths {
		if p == "stdout" {
			writers = append(writers, os.Stdout)
		} else {
			// Open or create file for logging.
			f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				writers = append(writers, f)
			} else {
				// Fallback to stdout if file cannot be opened.
				writers = append(writers, os.Stdout)
			}
		}
	}
	return io.MultiWriter(writers...)
}
