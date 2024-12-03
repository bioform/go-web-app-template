package server

import (
	"net/http"

	"github.com/bioform/go-web-app-template/pkg/server/middleware"
	"github.com/bioform/go-web-app-template/pkg/server/session"
	"github.com/bioform/go-web-app-template/web_api/api"
	"github.com/bioform/go-web-app-template/web_api/authapi"
	"github.com/bioform/go-web-app-template/web_api/healthapi"
	"github.com/bioform/go-web-app-template/web_api/helloapi"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	api := api.New(mux)

	mux.Handle("GET /{$}", helloapi.HelloHandler())

	mux.Handle("/auth/", authapi.RegisterRoutes("/auth"))

	huma.Get(api, "/health", healthapi.HealthHandler)

	var handler http.Handler = mux

	handler = middleware.UserSession(handler)
	handler = middleware.ApiProvider(handler)
	handler = session.Manager.LoadAndSave(handler)
	handler = middleware.Error(handler)
	handler = middleware.Tracing(handler)

	return handler
}
