package helloapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bioform/go-web-app-template/pkg/database"
	"github.com/bioform/go-web-app-template/pkg/server/rest"
	"github.com/bioform/go-web-app-template/pkg/server/session"
	"gorm.io/gorm"
)

var serveError = rest.ServeError

// HelloHandler returns an HTTP handler function that responds with a "Hello, world" message.
func HelloHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serveError(newHelloHandler(r.Context()).handle).ServeHTTP(w, r) // pass the "app" context copy to the handler
	}
}

type helloHandler struct {
	db *gorm.DB
}

func newHelloHandler(ctx context.Context) *helloHandler {
	return &helloHandler{
		db: database.Get(ctx),
	}
}

// hello responds to the request with a plain-text "Hello, world" message.
func (h *helloHandler) handle(w http.ResponseWriter, r *http.Request) error {
	buf := bytes.NewBufferString("")

	host, _ := os.Hostname()
	user := session.GetUser(r.Context())

	if user == nil {
		fmt.Fprintf(buf, "Hello,world!\n")
	} else {
		fmt.Fprintf(buf, "Hello, %s!\n", user.Name)
	}

	fmt.Fprintf(buf, "Version: 1.0.0\n")
	fmt.Fprintf(buf, "Hostname: %s\n", host)
	fmt.Fprintf(buf, "Database: %s\n", h.db.Name())

	resp := make(map[string]string)
	resp["message"] = buf.String()

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)

	return nil
}
