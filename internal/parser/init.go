package parser

import (
	"fmt"
	"log/slog"
)

func (p *Parser) InitializeStoreID() error {
	// Приоритет у координат
	if p.config.Longitude != 0.0 && p.config.Latitude != 0.0 {
		storeID, err := p.GetStoreIDByCoords(p.config.Longitude, p.config.Latitude)

		if err != nil {
			slog.Error("Couldn't identify store by coordinates", "error", err)
			return fmt.Errorf("failed to detect store by coords: %w", err)
		}

		p.client.SetStoreID(storeID)
		return nil
	}

	slog.Info("Using default store ID from config", "store_id", p.config.StoreID)
	p.client.SetStoreID(p.config.StoreID)
	return nil
}
