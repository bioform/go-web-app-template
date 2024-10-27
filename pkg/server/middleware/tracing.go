package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/bioform/go-web-app-template/pkg/logging"
)

func Tracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ctx := logging.AssignTraceID(r.Context())
		trace_id := logging.GetTraceID(ctx)

		r = r.WithContext(ctx)
		defer func() {
			slog.Debug("Completed request", "time", time.Since(start), "trace_id", trace_id)
		}()

		slog.Debug("Incomming request", "method", r.Method, "path", r.RequestURI, "remote_addr", r.RemoteAddr, "trace_id", trace_id)

		next.ServeHTTP(w, r)
	})
}
