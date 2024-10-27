package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bioform/go-web-app-template/pkg/env"
)

func NewHttpServer() *http.Server {
	port, err := strconv.Atoi(env.Get("PORT", "8080"))
	if err != nil {
		log.Panicf("error parsing PORT env var: %v", err)
	}

	// Declare Server config
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return httpServer
}
