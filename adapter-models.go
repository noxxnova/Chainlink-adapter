package main

// ChainlinkRequest represents the structure of incoming Chainlink job requests
type ChainlinkRequest struct {
	ID   string      `json:"id"`
	Data RequestData `json:"data"`
}

// RequestData contains the cryptocurrency symbol to query
type RequestData struct {
	Symbol string `json:"symbol"`
}

// ChainlinkResponse represents the structure of outgoing Chainlink job responses
type ChainlinkResponse struct {
	JobRunID     string      `json:"jobRunID"`
	StatusCode   int         `json:"statusCode"`
	Status       string      `json:"status"`
	Data         ResponseData `json:"data"`
	Error        string      `json:"error,omitempty"`
	Pending      bool        `json:"pending,omitempty"`
}

// ResponseData contains the price data retrieved from the exchange
type ResponseData struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Source string  `json:"source"`
}

// Error response helpers
func errorResponse(jobRunID string, message string, statusCode int) ChainlinkResponse {
	return ChainlinkResponse{
		JobRunID:   jobRunID,
		StatusCode: statusCode,
		Status:     "errored",
		Error:      message,
	}
}

// Success response helper
func successResponse(jobRunID string, data ResponseData) ChainlinkResponse {
	return ChainlinkResponse{
		JobRunID:   jobRunID,
		StatusCode: 200,
		Status:     "success",
		Data:       data,
	}
}