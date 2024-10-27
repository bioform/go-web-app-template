package authapi

import (
	"net/http"

	"github.com/bioform/go-web-app-template/web_api/authapi/handler"
)

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /login", handler.LoginHandler())

	return mux
}
