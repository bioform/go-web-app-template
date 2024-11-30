package logging

import (
	"context"
	"log/slog"

	"github.com/bioform/go-web-app-template/pkg/util/ctxstore"
)

func Logger(ctx context.Context) *slog.Logger {
	return ChildLogger(ctx, slog.Default())
}

func ChildLogger(ctx context.Context, ancestorLogger *slog.Logger) *slog.Logger {
	logger := ancestorLogger
	if logger == nil {
		logger = slog.Default()
	}

	traceID := ctxstore.GetTraceID(ctx)
	user := ctxstore.GetUser(ctx)

	if traceID != "" {
		logger = logger.With(slog.String("trace_id", traceID))
	}
	if user != nil {
		logger = logger.With(slog.Any("user", user))
	}

	return logger
}
