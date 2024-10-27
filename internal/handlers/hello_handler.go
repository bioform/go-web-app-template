package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bioform/go-web-app-template/internal/database"
)

// hello responds to the request with a plain-text "Hello, world" message.
func HandleHello(s *database.DbRefStorer, w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)

	buf := bytes.NewBufferString("")

	host, _ := os.Hostname()
	fmt.Fprintf(buf, "Hello, world!\n")
	fmt.Fprintf(buf, "Version: 1.0.0\n")
	fmt.Fprintf(buf, "Hostname: %s\n", host)

	resp := make(map[string]string)
	resp["message"] = buf.String()

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
