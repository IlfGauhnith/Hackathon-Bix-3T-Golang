package model

// CSVRecord represents one line in the uploaded CSV.
// Field tags match the CSV headers exactly.
type CSVRecord struct {
	ID       int     `csv:"id"`         // Unique product ID
	Name     string  `csv:"nome"`       // Product name
	Category string  `csv:"categoria"`  // Product category
	Price    float64 `csv:"preco"`      // Product price
	Stock    int     `csv:"estoque"`    // Product stock quantity
	Supplier string  `csv:"fornecedor"` // Product supplier name
}
