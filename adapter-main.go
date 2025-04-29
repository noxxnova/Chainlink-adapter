package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Constants
const (
	DefaultPort            = "8080"
	DefaultTimeout         = 10 * time.Second
	ContentTypeJSON        = "application/json"
	RequestTimeout         = 5 * time.Second
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	// Setup HTTP server with routes
	http.HandleFunc("/coingecko", handleCoinbaseRequest)
	http.HandleFunc("/binance", handleBinanceRequest)
	http.HandleFunc("/okx", handleOKXRequest)
	http.HandleFunc("/health", handleHealthCheck)

	// Start the server
	fmt.Printf("Starting Chainlink External Adapter server on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Handler for health check endpoint
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
