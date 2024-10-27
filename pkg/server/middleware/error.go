package middleware

import (
	"fmt"
	"net/http"

	"golang.org/x/exp/slog"
)

func Error(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Recover from panics and handle errors
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf("Recovered from panic: %v", err)
				slog.Error(msg, "error", err)

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
