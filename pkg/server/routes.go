package server

import (
	"net/http"

	"github.com/bioform/go-web-app-template/pkg/server/middleware"
	"github.com/bioform/go-web-app-template/web_api/authapi"
	"github.com/bioform/go-web-app-template/web_api/healthapi"
	"github.com/bioform/go-web-app-template/web_api/helloapi"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

func RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	api := humago.New(mux, huma.DefaultConfig("My API", "1.0.0"))

	mux.Handle("GET /{$}", helloapi.HelloHandler())

	mux.Handle("/auth/", http.StripPrefix("/auth", authapi.RegisterRoutes()))

	huma.Get(api, "/health", healthapi.HealthHandler)

	var handler http.Handler = mux

	handler = middleware.Error(handler)
	handler = middleware.Tracing(handler)

	return handler
}
