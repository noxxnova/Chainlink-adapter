package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Common request handling function
func processRequest(w http.ResponseWriter, r *http.Request, fetchPrice func(string) (float64, error), source string) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var chainlinkReq ChainlinkRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendErrorResponse(w, "", "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &chainlinkReq)
	if err != nil {
		sendErrorResponse(w, "", "Error parsing request JSON", http.StatusBadRequest)
		return
	}

	// Get the symbol from the request
	symbol := chainlinkReq.Data.Symbol
	if symbol == "" {
		sendErrorResponse(w, chainlinkReq.ID, "Symbol is required", http.StatusBadRequest)
		return
	}

	// Fetch the price from the exchange
	price, err := fetchPrice(symbol)
	if err != nil {
		sendErrorResponse(w, chainlinkReq.ID, fmt.Sprintf("Error fetching price: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Create the response
	responseData := ResponseData{
		Symbol: symbol,
		Price:  price,
		Source: source,
	}

	response := successResponse(chainlinkReq.ID, responseData)

	// Send the response
	w.Header().Set("Content-Type", ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Helper function to send error responses
func sendErrorResponse(w http.ResponseWriter, jobRunID string, message string, statusCode int) {
	w.Header().Set("Content-Type", ContentTypeJSON)
	w.WriteHeader(statusCode)
	response := errorResponse(jobRunID, message, statusCode)
	json.NewEncoder(w).Encode(response)
}

// Exchange-specific request handlers
func handleCoinbaseRequest(w http.ResponseWriter, r *http.Request) {
	processRequest(w, r, fetchCoinbasePrice, "CoinGecko")
}

func handleBinanceRequest(w http.ResponseWriter, r *http.Request) {
	processRequest(w, r, fetchBinancePrice, "Binance")
}

func handleOKXRequest(w http.ResponseWriter, r *http.Request) {
	processRequest(w, r, fetchOKXPrice, "OKX")
}
