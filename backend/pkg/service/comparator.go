package service

import (
	"github.com/IlfGauhnith/Hackathon-Bix-3T-Golang/pkg/model"
)

// CompareBatch fetches one page of external products and compares them to the given CSV slice.
func CompareBatch(records []model.CSVRecord, apiResp *model.APIResponse) ([]model.Divergence, error) {

	// Build a lookup map by product ID
	externalMap := make(map[int]model.ExternalProduct, len(apiResp.Data))
	for _, p := range apiResp.Data {
		externalMap[p.ID] = p
	}

	var results []model.Divergence

	// Compare each CSV record
	for _, rec := range records {
		ext, found := externalMap[rec.ID]
		if !found {
			// Entire record missing in API
			results = append(results, model.Divergence{
				RecordID: rec.ID,
				Differences: []model.FieldDifference{{
					FieldName: "existence",
					CSVValue:  true,
					APIValue:  false,
				}},
			})
			continue
		}

		var diffs []model.FieldDifference

		if rec.Name != ext.Name {
			diffs = append(diffs, model.FieldDifference{"Name", rec.Name, ext.Name})
		}
		if rec.Category != ext.Category {
			diffs = append(diffs, model.FieldDifference{"Category", rec.Category, ext.Category})
		}
		if rec.Price != ext.Price {
			diffs = append(diffs, model.FieldDifference{"Price", rec.Price, ext.Price})
		}
		if rec.Stock != ext.Stock {
			diffs = append(diffs, model.FieldDifference{"Stock", rec.Stock, ext.Stock})
		}
		if rec.Supplier != ext.Supplier {
			diffs = append(diffs, model.FieldDifference{"Supplier", rec.Supplier, ext.Supplier})
		}

		if len(diffs) > 0 {
			results = append(results, model.Divergence{RecordID: rec.ID, Differences: diffs})
		}
	}

	return results, nil
}
