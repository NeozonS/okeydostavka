package parser

import (
	"TestCenozavr/internal/models"
	"encoding/json"
	"log/slog"
)

func (p *Parser) GetCatalog() (models.Categories, error) {
	resp, err := p.client.GetCategory()
	if err != nil {
		slog.Error("Error receiving the catalog", "error", err)
		return models.Categories{}, err
	}
	var fullCatalog models.Categories

	if err := json.Unmarshal(resp.BodyBytes, &fullCatalog); err != nil {
		slog.Error("Catalog parsing error", "error", err, "body_length", len(resp.BodyBytes))
		return models.Categories{}, err
	}
	categories := searchLeafCategories(fullCatalog.Categories)

	return models.Categories{Categories: categories}, nil
}

func searchLeafCategories(cat []models.Category) []models.Category {
	var result []models.Category

	for _, c := range cat {
		if len(c.Subcategories) == 0 {
			result = append(result, models.Category{
				ID:            c.ID,
				Name:          c.Name,
				Subcategories: nil,
			})
		} else {
			result = append(result, searchLeafCategories(c.Subcategories)...)
		}
	}
	return result
}
