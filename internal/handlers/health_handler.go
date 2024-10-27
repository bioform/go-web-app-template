package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bioform/go-web-app-template/internal/database"
)

func HealthHandler(s *database.DbRefStorer, w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.Db.Health())

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
