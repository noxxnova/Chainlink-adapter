package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// createHTTPClient returns an HTTP client with timeout
func createHTTPClient() *http.Client {
	return &http.Client{
		Timeout: RequestTimeout,
	}
}

// CoinGecko API client
func fetchCoinbasePrice(symbol string) (float64, error) {
	// Map common symbols to CoinGecko IDs
	symbolMap := map[string]string{
		"BTC":  "bitcoin",
		"ETH":  "ethereum",
		"LINK": "chainlink",
		// Add more mappings as needed
	}

	// Convert symbol to lowercase for comparison
	upperSymbol := strings.ToUpper(symbol)
	coinID, exists := symbolMap[upperSymbol]
	if !exists {
		return 0, fmt.Errorf("unsupported symbol: %s", symbol)
	}

	// Create the request
	client := createHTTPClient()
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", coinID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bad response from CoinGecko: %d", resp.StatusCode)
	}

	// Parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %w", err)
	}

	// Unmarshal dynamic JSON response
	var response map[string]map[string]float64
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, fmt.Errorf("error parsing response JSON: %w", err)
	}

	// Access the price using the coin ID
	coinData, exists := response[coinID]
	if !exists {
		return 0, fmt.Errorf("coin data not found in response")
	}

	price, exists := coinData["usd"]
	if !exists {
		return 0, fmt.Errorf("USD price not found in response")
	}

	return price, nil
}

// Binance API client
func fetchBinancePrice(symbol string) (float64, error) {
	// Format the symbol for Binance API (BTCUSDT, ETHUSDT, etc.)
	formattedSymbol := fmt.Sprintf("%sUSDT", strings.ToUpper(symbol))

	// Create the request
	client := createHTTPClient()
	url := "https://api.binance.com/api/v3/ticker/price"

	// If symbol is provided, add it as a query parameter
	if symbol != "" {
		url = fmt.Sprintf("%s?symbol=%s", url, formattedSymbol)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bad response from Binance: %d", resp.StatusCode)
	}

	// Parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %w", err)
	}

	// Handle both single item response and array response
	if strings.HasPrefix(string(body), "[") {
		// Array response (when no symbol is specified)
		var responseArray []struct {
			Symbol string `json:"symbol"`
			Price  string `json:"price"`
		}

		err = json.Unmarshal(body, &responseArray)
		if err != nil {
			return 0, fmt.Errorf("error parsing response JSON array: %w", err)
		}

		// Find the matching symbol
		for _, item := range responseArray {
			if item.Symbol == formattedSymbol {
				var price float64
				_, err = fmt.Sscanf(item.Price, "%f", &price)
				if err != nil {
					return 0, fmt.Errorf("error parsing price: %w", err)
				}
				return price, nil
			}
		}

		return 0, fmt.Errorf("symbol %s not found in Binance response", formattedSymbol)
	} else {
		// Single item response (when symbol is specified)
		var response struct {
			Symbol string `json:"symbol"`
			Price  string `json:"price"`
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return 0, fmt.Errorf("error parsing response JSON: %w", err)
		}

		// Parse the price as float
		var price float64
		_, err = fmt.Sscanf(response.Price, "%f", &price)
		if err != nil {
			return 0, fmt.Errorf("error parsing price: %w", err)
		}

		return price, nil
	}
}

// OKX API client
func fetchOKXPrice(symbol string) (float64, error) {
	// Format the symbol for OKX API (BTC-USDT, ETH-USDT, etc.)
	formattedSymbol := fmt.Sprintf("%s-USDT", strings.ToUpper(symbol))

	// Create the request
	client := createHTTPClient()
	url := "https://www.okx.com/api/v5/market/tickers?instType=SPOT"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bad response from OKX: %d", resp.StatusCode)
	}

	// Parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %w", err)
	}

	var response struct {
		Data []struct {
			InstID string `json:"instId"`
			Last   string `json:"last"`
			LastPx string `json:"lastPx,omitempty"`
			AskPx  string `json:"askPx,omitempty"`
			BidPx  string `json:"bidPx,omitempty"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, fmt.Errorf("error parsing response JSON: %w", err)
	}

	// Find the ticker with matching symbol
	for _, ticker := range response.Data {
		if ticker.InstID == formattedSymbol {
			// Determine which price field to use (last is preferred)
			priceStr := ticker.Last
			if priceStr == "" {
				priceStr = ticker.LastPx
			}
			if priceStr == "" {
				// If no last price, use mid price between bid and ask
				var askPrice, bidPrice float64
				if ticker.AskPx != "" && ticker.BidPx != "" {
					fmt.Sscanf(ticker.AskPx, "%f", &askPrice)
					fmt.Sscanf(ticker.BidPx, "%f", &bidPrice)
					return (askPrice + bidPrice) / 2, nil
				}
				return 0, fmt.Errorf("no price data for %s", formattedSymbol)
			}

			// Parse the price as float
			var price float64
			_, err = fmt.Sscanf(priceStr, "%f", &price)
			if err != nil {
				return 0, fmt.Errorf("error parsing price: %w", err)
			}

			return price, nil
		}
	}

	return 0, fmt.Errorf("symbol %s not found in OKX response", formattedSymbol)
}
