package model

// HealthResponse holds status for the health check.
// swagger:model
type HealthResponse struct {
	// API status message
	// example: API is running
	Status string `json:"status"`
}
