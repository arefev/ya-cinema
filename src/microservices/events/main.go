package main

import (
	"log"
	"net/http"
	"os"
)


func main() {
	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082" // Note: Using a different port than the monolith
	}
	log.Printf("Starting events microservice on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
