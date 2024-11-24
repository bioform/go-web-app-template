package logging

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// ContextKey is used for context.Context value. The value requires a key that is not primitive type.
type ContextKey string // can be unexported

// ContextKeyTraceID is the ContextKey for TraceID
const ContextKeyTraceID ContextKey = "trace_id" // can be unexported

func Get(ctx context.Context) *slog.Logger {
	return Child(ctx, slog.Default())
}

func Child(ctx context.Context, ancestorLogger *slog.Logger) *slog.Logger {

	traceID := GetTraceID(ctx)

	// initiate the variable that holds the logger being returned at the end, make standard logger the default
	logger := slog.Default()

	if ancestorLogger != nil {
		logger = ancestorLogger
	}

	if len(traceID) == 0 {
		return logger
	}

	logger = logger.With(slog.String(string(ContextKeyTraceID), traceID))

	return logger
}

// AttachTraceID will attach a brand new request ID to a http request
func AssignTraceID(ctx context.Context) context.Context {

	reqID := uuid.NewString()

	return context.WithValue(ctx, ContextKeyTraceID, reqID)
}

// GetTraceID will get reqID from a http request and return it as a string
func GetTraceID(ctx context.Context) string {

	traceIDValue := ctx.Value(ContextKeyTraceID)

	if traceID, ok := traceIDValue.(string); ok {
		return traceID
	}

	return ""
}
