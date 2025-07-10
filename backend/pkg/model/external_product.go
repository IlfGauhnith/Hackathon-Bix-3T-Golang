package model

// ExternalProduct mirrors the JSON structure returned by the API.
type ExternalProduct struct {
	ID       int     `json:"id"`
	Name     string  `json:"nome"`
	Category string  `json:"categoria"`
	Price    float64 `json:"preco"`
	Stock    int     `json:"estoque"`
	Supplier string  `json:"fornecedor"`
}

// Pagination holds pagination metadata from the external API.
type Pagination struct {
	CurrentPage     int  `json:"current_page"`
	ItemsPerPage    int  `json:"items_per_page"`
	TotalItems      int  `json:"total_items"`
	TotalPages      int  `json:"total_pages"`
	HasNextPage     bool `json:"has_next_page"`
	HasPreviousPage bool `json:"has_previous_page"`
}

// APIResponse models the full JSON payload from the external API.
type APIResponse struct {
	Data       []ExternalProduct `json:"data"`
	Pagination Pagination        `json:"pagination"`
}
