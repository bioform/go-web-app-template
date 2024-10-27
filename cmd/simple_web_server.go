package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

// hello responds to the request with a plain-text "Hello, world" message.
func handleHello(w http.ResponseWriter, r *http.Request) {
	host, _ := os.Hostname()

	fmt.Fprintf(w, "Hello, world!\n")
	fmt.Fprintf(w, "Version: 1.0.0\n")
	fmt.Fprintf(w, "Hostname: %s\n", host)
}

func main() {
	listenAddr := flag.String("listenaddr", ":8080", "http service address")
	flag.Parse()

	// register hello function to handle all requests
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", http.HandlerFunc(handleHello))

	// start the web server on port and accept requests
	log.Printf("Server listening on port %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, mux))
}
