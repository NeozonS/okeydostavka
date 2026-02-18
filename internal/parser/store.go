package parser

import (
	"TestCenozavr/internal/models"
	"encoding/json"
	"fmt"
	"log/slog"
)

func (p *Parser) GetStoreIDByCoords(lon, lat float64) (uint64, error) {
	resp, err := p.client.GetStoreByCoords(lon, lat)
	if err != nil {
		return 0, err
	}

	var store models.AddressResponse

	if err := json.Unmarshal(resp.BodyBytes, &store); err != nil {
		slog.Error("Coordinate parsing error", "error", err, "latitude", lat, "longitude", lon, "body_length", len(resp.BodyBytes))
		return 0, fmt.Errorf("failed to unmarshal store response: %w", err)
	}

	if len(store.Addresses) == 0 {
		slog.Warn("No stores found for coordinates", "latitude", lat, "longitude", lon)
		return 0, fmt.Errorf("the store does not deliver at these coordinates, or the coordinates are specified incorrectly (%f, %f)", lat, lon)
	}

	result := uint64(store.Addresses[0].FfcID)

	return result, nil
}
