package middleware

import (
	"net/http"

	"github.com/bioform/go-web-app-template/pkg/api"
	"github.com/bioform/go-web-app-template/pkg/database"
)

func ApiProvider(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := database.Default()
		api := api.New(db)
		ctx := api.AddTo(r.Context())
		// Call the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
