package logx

import "context"

// Empty struct types are used as context keys to avoid collisions.
type (
	traceIDKey struct{}
	jobIDKey   struct{}
)

// WithTraceID returns a new context that contains the given traceID.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

// TraceIDFromContext extracts the traceID from the context if available.
func TraceIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(traceIDKey{}).(string)
	return v, ok
}

// WithJobID returns a new context that contains the given jobID.
func WithJobID(ctx context.Context, jobID int64) context.Context {
	return context.WithValue(ctx, jobIDKey{}, jobID)
}

// JobIDFromContext extracts the jobID from the context if available.
func JobIDFromContext(ctx context.Context) (int64, bool) {
	v, ok := ctx.Value(jobIDKey{}).(int64)
	return v, ok
}
