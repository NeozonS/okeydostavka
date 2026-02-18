package client

import (
	"fmt"
	"log/slog"
	"net/url"

	"github.com/Danny-Dasilva/CycleTLS/cycletls"
)

func (c *ClientWithCookies) GetCategory() (cycletls.Response, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		slog.Error("Failed to parse base URL", "error", err)
		return cycletls.Response{}, fmt.Errorf("failed to parse base URL: %w", err)
	}
	u = u.JoinPath("wcs", "resources", "mobihub023", "store", fmt.Sprintf("%d", c.StoreID), "catalog", "categories")

	slog.Debug("Requesting category",
		"url", u.String(),
		"store_id", c.StoreID,
	)

	response, err := c.Do(u.String(), cycletls.Options{
		Ja3:       c.ja3,
		UserAgent: c.UA,
		Proxy:     c.Proxy,

		Headers: map[string]string{
			"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			"Accept-Language":           "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7",
			"Accept-Encoding":           "gzip, deflate, br, zstd",
			"Upgrade-Insecure-Requests": "1",
			"Sec-Fetch-Dest":            "document",
			"Sec-Fetch-Mode":            "navigate",
			"Sec-Fetch-Site":            "none",
			"Priority":                  "u=0, i",
			"TE":                        "trailers",
			"Connection":                "keep-alive",
		},
		Timeout:               30,
		EnableConnectionReuse: true,
	}, "GET")

	if err != nil {
		slog.Error("Error in the catalog request", "error", err, "store_id", c.StoreID)
		return cycletls.Response{}, err
	}
	if response.Status != 200 {
		slog.Warn("Error receiving the Category",
			"status", response.Status,
			"store_id", c.StoreID,
		)
		return cycletls.Response{}, fmt.Errorf("unexpected status code: %d", response.Status)
	}

	slog.Info("Category received successfully", "store_id", c.StoreID, "status", response.Status)

	return response, err
}
