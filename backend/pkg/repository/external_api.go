package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/IlfGauhnith/Hackathon-Bix-3T-Golang/pkg/config"
	"github.com/IlfGauhnith/Hackathon-Bix-3T-Golang/pkg/model"
)

// ExternalAPIClient handles HTTP calls to the product API.
type ExternalAPIClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewExternalAPIClient constructs a client with the given config.
func NewExternalAPIClient(cfg *config.Config) *ExternalAPIClient {
	return &ExternalAPIClient{
		baseURL:    cfg.ExternalAPIBaseURL,
		httpClient: &http.Client{},
	}
}

// GetProducts fetches a single page of products.
func (c *ExternalAPIClient) GetProducts(page, limit int) (*model.APIResponse, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	q := u.Query()
	q.Set("page", fmt.Sprintf("%d", page))
	q.Set("limit", fmt.Sprintf("%d", limit))
	u.RawQuery = q.Encode()

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("external API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var apiResp model.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}
	return &apiResp, nil
}
