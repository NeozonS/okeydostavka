package client

import (
	"fmt"
	"log/slog"
	"net/url"

	"github.com/Danny-Dasilva/CycleTLS/cycletls"
)

func (c *ClientWithCookies) GetCity() (cycletls.Response, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		slog.Error("Failed to parse base URL", "error", err)
		return cycletls.Response{}, fmt.Errorf("failed to parse base URL: %w", err)
	}

	u = u.JoinPath("wcs", "resources", "mobihub023", "store", fmt.Sprintf("%d", c.StoreID), "regions", "online")

	response, err := c.Do(u.String(), cycletls.Options{
		Body:      "",
		Ja3:       c.ja3,
		UserAgent: c.UA,
		Proxy:     c.Proxy,

		Headers: map[string]string{
			"Accept":           "application/json, text/plain, */*",
			"Content-Type":     "application/json",
			"Origin":           baseURL,
			"X-Requested-With": "XMLHttpRequest",
		},
		Timeout:               30,
		EnableConnectionReuse: true,
	},
		"GET",
	)
	if err != nil {
		slog.Error("Error requesting a city", "error", err, "store_id", c.StoreID)
		return cycletls.Response{}, err
	}
	if response.Status != 200 {
		slog.Warn("Error receiving the City", "status", response.Status, "store_id", c.StoreID)
		return cycletls.Response{}, fmt.Errorf("unexpected status code: %d", response.Status)
	}

	slog.Info("City received successfully", "store_id", c.StoreID, "status", response.Status)

	return response, err
}
