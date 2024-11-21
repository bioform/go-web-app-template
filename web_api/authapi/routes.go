package authapi

import (
	"net/http"

	"github.com/bioform/go-web-app-template/web_api/api"
	"github.com/bioform/go-web-app-template/web_api/authapi/handler"
	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(prefix string) *http.ServeMux {
	mux := http.NewServeMux()
	api := api.NewWithPrefix(mux, prefix)

	huma.Post(api, "/login", handler.LoginHandler)

	return mux
}
