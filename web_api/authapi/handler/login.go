package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bioform/go-web-app-template/internal/user/repository"
	"github.com/bioform/go-web-app-template/pkg/server/rest"
)

func LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rest.ServeError(newHelloHandler(r.Context()).Handle).ServeHTTP(w, r) // pass the "app" context copy to the handler
	}
}

type loginHandler struct {
	repo repository.UserRepository
}

func newHelloHandler(ctx context.Context) *loginHandler {
	return &loginHandler{
		repo: repository.NewUserRepository(ctx),
	}
}

// hello responds to the request with a plain-text "Hello, world" message.
func (h *loginHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	resp := make(map[string]string)
	resp["message"] = "Success Login"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)

	return nil
}
