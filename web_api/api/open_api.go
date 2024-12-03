package api

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

var config = huma.DefaultConfig("My API", "1.0.0")

func init() {
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearer": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}
	config.Servers = []*huma.Server{{URL: "/api"}}
}

func New(mux humago.Mux) huma.API {
	return humago.New(mux, config)
}

func NewWithPrefix(mux humago.Mux, prefix string) huma.API {
	return humago.NewWithPrefix(mux, prefix, config)
}
