package model

// FieldDifference describes a single mismatched field.
type FieldDifference struct {
	FieldName string      `json:"field_name"` // Name of the differing field
	CSVValue  interface{} `json:"csv_value"`  // Value found in CSV
	APIValue  interface{} `json:"api_value"`  // Value returned by API
}

// Divergence aggregates all differences for one record ID.
type Divergence struct {
	RecordID    int               `json:"record_id"`   // ID of the CSV record
	Differences []FieldDifference `json:"differences"` // List of mismatches
}
