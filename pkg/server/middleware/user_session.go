package middleware

import (
	"context"
	"net/http"

	"github.com/bioform/go-web-app-template/internal/user/repository"
	"github.com/bioform/go-web-app-template/pkg/logging"
	"github.com/bioform/go-web-app-template/pkg/server/session"
)

func UserSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := session.Manager.GetInt64(ctx, session.UserIdKey)

		if userID != 0 {
			logger := logging.Logger(ctx)
			repo := repository.NewUserRepository()
			user, err := repo.FindByID(ctx, uint(userID))

			if err == nil {
				// Add the user to the request context
				ctx := context.WithValue(ctx, session.UserKey, user)
				r = r.WithContext(ctx)

				logger.Info("User loaded from session", "user", user)
			} else {
				// Clear the session if the user can't be found
				session.Manager.Remove(ctx, session.UserIdKey)
				logger.Error("Failed to load user from session", "error", err)
			}
		}

		next.ServeHTTP(w, r)
	})
}
