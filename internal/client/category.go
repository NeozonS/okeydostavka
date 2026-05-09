package client

import (
	"fmt"
	"log/slog"
	"net/url"
)

func (c *ClientWithCookies) GetCategory() ([]byte, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		slog.Error("Failed to parse base URL", "error", err)
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	u = u.JoinPath("wcs", "resources", "mobihub023", "store", fmt.Sprintf("%d", c.StoreID), "catalog", "categories")

	slog.Debug("Requesting category", "url", u.String(), "store_id", c.StoreID)

	resp := c.GET(u.String()).
		SetHeaders("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7").
		AddHeaders("Referer", baseURL+"/").
		AddHeaders("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8").
		AddHeaders("Upgrade-Insecure-Requests", "1").
		AddHeaders("Sec-Fetch-Dest", "document").
		AddHeaders("Sec-Fetch-Mode", "navigate").
		AddHeaders("Sec-Fetch-Site", "none").
		AddHeaders("Priority", "u=0, i").
		AddHeaders("TE", "trailers").
		Do()

	if resp.IsErr() {
		slog.Error("Error in the catalog request", "error", resp.Err(), "store_id", c.StoreID)
		return nil, resp.Err()
	}

	r := resp.Ok()
	if r.StatusCode != 200 {
		slog.Warn("Error receiving the Category", "status", r.StatusCode, "store_id", c.StoreID)
		r.Debug().Request().Response(true).Print()
		return nil, fmt.Errorf("unexpected status code: %d", r.StatusCode)
	}

	body := r.Body.Bytes()
	if body.IsErr() {
		return nil, body.Err()
	}

	slog.Info("Category received successfully", "store_id", c.StoreID, "status", r.StatusCode)
	return body.Ok(), nil
}
