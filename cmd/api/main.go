package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bioform/go-web-app-template/pkg/database"
	"github.com/bioform/go-web-app-template/pkg/server"
)

func main() {

	httpServer := server.NewHttpServer()

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go server.GracefulShutdown(httpServer, done)

	// start the web server on port and accept requests
	log.Printf("Server listening on port %s", httpServer.Addr)
	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	defer database.Close()

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
