package middleware

import (
	"net/http"

	"github.com/bioform/go-web-app-template/internal/user/repository"
	"github.com/bioform/go-web-app-template/pkg/api"
	"github.com/bioform/go-web-app-template/pkg/logging"
	"github.com/bioform/go-web-app-template/pkg/server/session"
	"github.com/bioform/go-web-app-template/pkg/util/ctxstore"
)

func UserSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := session.Manager.GetInt64(ctx, session.UserIdKey)

		if userID != 0 {
			logger := logging.Logger(ctx)

			api, err := api.From(ctx)
			if err != nil {
				logger.Error("user session middleware", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			db := api.DB()
			repo := repository.NewUserRepository(db)
			user, err := repo.FindByID(ctx, uint(userID))

			if err == nil {
				// Add the user to the request context
				ctx := ctxstore.AssignUser(ctx, user)
				r = r.WithContext(ctx)

				logging.Logger(ctx).Debug("User loaded from session")
			} else {
				// Clear the session if the user can't be found
				session.Manager.Remove(ctx, session.UserIdKey)
				logger.Error("user session middleware", "error", err)
			}
		}

		next.ServeHTTP(w, r)
	})
}
