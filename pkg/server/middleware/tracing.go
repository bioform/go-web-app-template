package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/bioform/go-web-app-template/pkg/logging"
	"github.com/bioform/go-web-app-template/pkg/util/ctxstore"
)

func Tracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Recover from panics and handle errors
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf("tracing middleware panic: %v", err)
				slog.Error(msg, "error", err)

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		start := time.Now()

		ctx := ctxstore.AssignTraceID(r.Context())
		log := logging.Logger(ctx)

		r = r.WithContext(ctx)
		defer func() {
			log.Debug("Completed request", "time", time.Since(start))
		}()

		log.Debug("Incomming request", "method", r.Method, "path", r.RequestURI, "remote_addr", r.RemoteAddr)

		next.ServeHTTP(w, r)
	})
}
