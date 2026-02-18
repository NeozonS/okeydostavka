package parser

import (
	"TestCenozavr/internal/models"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

func (p *Parser) GetCity() (models.Regions, error) {
	var out models.Regions
	resp, err := p.client.GetCity()
	if err != nil {
		slog.Error("Error getting cities", "error", err)
		return models.Regions{}, err
	}

	if err := json.Unmarshal(resp.BodyBytes, &out); err != nil {
		slog.Error("City parsing error", "error", err, "body_length", len(resp.BodyBytes))
		return models.Regions{}, fmt.Errorf("failed to unmarshal cities: %w", err)
	}
	return out, nil
}

func (p *Parser) GetStoreForAllCities() (models.Regions, error) {
	city, err := p.GetCity()
	if err != nil {
		slog.Error("Failed to get cities for store lookup", "error", err)
		return models.Regions{}, err
	}

	for i, _ := range city.Regions {

		time.Sleep(1 * time.Second)

		respstore, err := p.client.GetStore(city.Regions[i].FfcID)

		if err != nil {
			slog.Warn("Error getting stores for city", "city", city.Regions[i].City, "ffc_id", city.Regions[i].FfcID, "error", err)
			continue
		}

		var store models.Stores
		if err := json.Unmarshal(respstore.BodyBytes, &store); err != nil {
			slog.Warn("Error parsing stores for city", "city", city.Regions[i].City, "ffc_id", city.Regions[i].FfcID, "error", err)
			continue
		}
		city.Regions[i].Stores = store.Stores
	}
	return city, nil
}
