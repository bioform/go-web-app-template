package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/bioform/go-web-app-template/internal/database"
	"github.com/bioform/go-web-app-template/internal/utils"
)

type Server struct {
	*database.DbRefStorer

	port int
}

func New() *Server {
	port, _ := strconv.Atoi(utils.GetEnv("PORT", "8080"))
	return &Server{
		port: port,

		DbRefStorer: &database.DbRefStorer{
			Db: database.New(),
		},
	}
}

func NewHttpServer() *http.Server {
	server := New()

	// Declare Server config
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", server.port),
		Handler:      server.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return httpServer
}
