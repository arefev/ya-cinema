package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)


func main() {
	http.HandleFunc("/api/events/health", handleHealth)
	
	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082" // Note: Using a different port than the monolith
	}
	log.Printf("Starting events microservice on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"status": true})
}