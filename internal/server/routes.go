package server

import (
	"net/http"

	"github.com/bioform/go-web-app-template/internal/database"
	"github.com/bioform/go-web-app-template/internal/handlers"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.Handle("GET /{$}", s.newHandler(handlers.HandleHello)) // see: https://pkg.go.dev/net/http#ServeMux
	mux.Handle("GET /health", s.newHandler(handlers.HealthHandler))

	return mux
}

func (s *Server) newHandler(h func(*database.DbRefStorer, http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(s.DbRefStorer, w, r)
	})
}
