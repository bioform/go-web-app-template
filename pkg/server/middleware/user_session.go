package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/bioform/go-web-app-template/internal/user/repository"
	"github.com/bioform/go-web-app-template/pkg/server/session"
)

func UserSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := session.Manager.GetInt64(ctx, session.UserIdKey)

		if userID != 0 {
			repo := repository.NewUserRepository(ctx)
			user, err := repo.FindByID(uint(userID))
			if err == nil {
				// Add the user to the request context
				ctx := context.WithValue(ctx, session.UserKey, user)
				r = r.WithContext(ctx)
			} else {
				// Clear the session if the user can't be found
				session.Manager.Remove(ctx, session.UserIdKey)
				slog.Error("Failed to load user from session", "error", err)
			}
		}

		next.ServeHTTP(w, r)
	})
}