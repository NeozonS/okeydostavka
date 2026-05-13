package client

import (
	"fmt"
	"log/slog"
	"net/url"
)

func (c *Client) GetCity() ([]byte, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		slog.Error("Failed to parse base URL", "error", err)
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	u = u.JoinPath("wcs", "resources", "mobihub023", "store", fmt.Sprintf("%d", c.StoreID), "regions", "online")

	resp := c.GET(u.String()).
		AddHeaders("Referer", baseURL+"/").
		AddHeaders("Accept", "application/json, text/plain, */*").
		AddHeaders("Content-Type", "application/json").
		AddHeaders("Origin", baseURL).
		AddHeaders("X-Requested-With", "XMLHttpRequest").
		Do()

	if resp.IsErr() {
		slog.Error("Error requesting a city", "error", resp.Err(), "store_id", c.StoreID)
		return nil, resp.Err()
	}

	r := resp.Ok()
	if r.StatusCode != 200 {
		slog.Warn("Error receiving the City", "status", r.StatusCode, "store_id", c.StoreID)
		return nil, fmt.Errorf("unexpected status code: %d", r.StatusCode)
	}

	body := r.Body.Bytes()
	if body.IsErr() {
		return nil, body.Err()
	}

	slog.Info("City received successfully", "store_id", c.StoreID, "status", r.StatusCode)
	return body.Ok(), nil
}
