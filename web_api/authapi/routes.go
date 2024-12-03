package authapi

import (
	"github.com/bioform/go-web-app-template/web_api/api/route"
	"github.com/bioform/go-web-app-template/web_api/authapi/handler"
	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API, prefix string) {
	p := route.New(prefix).Path

	huma.Post(api, p("/login"), handler.LoginHandler)
}
