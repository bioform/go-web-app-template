package ctxstore

import (
	"context"
	"log/slog"

	"github.com/bioform/go-web-app-template/internal/user/model"
	"github.com/google/uuid"
)

// ContextKey is used for context.Context value. The value requires a key that is not primitive type.
type contextKey string // can be unexported

type userKey string // can be unexported

// TraceIDKey is the ContextKey for TraceID
const TraceIDKey contextKey = "trace_id" // can be unexported
// UserKey is the ContextKey for User
const UserKey userKey = "user"

// AttachTraceID will attach a brand new request ID to a http request
func AssignTraceID(ctx context.Context) context.Context {

	reqID := uuid.NewString()

	return context.WithValue(ctx, TraceIDKey, reqID)
}

// GetTraceID will get reqID from a http request and return it as a string
func GetTraceID(ctx context.Context) string {

	traceIDValue := ctx.Value(TraceIDKey)

	if traceID, ok := traceIDValue.(string); ok {
		return traceID
	}

	return ""
}

func AssignUser(ctx context.Context, user *model.User) context.Context {

	return context.WithValue(ctx, UserKey, user)
}

func GetUser(ctx context.Context) *model.User {
	u := ctx.Value(UserKey)
	if u == nil {
		return nil
	}

	user, ok := u.(*model.User)
	if !ok {
		traceID := GetTraceID(ctx)
		slog.With(slog.String("trace_id", traceID)).Error("get user from context", "user", u)
		return nil
	}

	return user
}
