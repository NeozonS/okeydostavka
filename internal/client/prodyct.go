package client

import (
	"fmt"
	"log/slog"
	"net/url"

	"github.com/Danny-Dasilva/CycleTLS/cycletls"
)

func (c *ClientWithCookies) GetProductInfo(catalogID string, pageNumber, pageSize int) (cycletls.Response, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		slog.Error("Failed to parse base URL", "error", err)
		return cycletls.Response{}, fmt.Errorf("failed to parse base URL: %w", err)
	}

	u = u.JoinPath("wcs", "resources", "mobihub023", "store", fmt.Sprintf("%d", c.StoreID), "catalog", "combo", catalogID)
	q := u.Query()
	q.Set("pageSize", fmt.Sprintf("%d", pageSize))
	q.Set("pageNumber", fmt.Sprintf("%d", pageNumber))
	u.RawQuery = q.Encode()

	response, err := c.Do(u.String(), cycletls.Options{
		Ja3:       c.ja3,
		UserAgent: c.UA,
		Proxy:     c.Proxy,

		Headers: map[string]string{
			"Accept":           "*/*",
			"Accept-Language":  "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7",
			"Content-Language": "ru-RU",
			"Content-Type":     "application/x-www-form-urlencoded",
			"Priority":         "u=4",
			"X-Requested-With": "XMLHttpRequest",
			"Sec-Fetch-Dest":   "empty",
			"Sec-Fetch-Mode":   "cors",
			"Sec-Fetch-Site":   "same-origin",
		},
		Timeout:               30,
		EnableConnectionReuse: true,
	},
		"POST",
	)
	if err != nil {
		slog.Error("Price receipt request error", "error", err, "store_id", c.StoreID, "catalog_id", catalogID)
		return cycletls.Response{}, err
	}
	if response.Status != 200 {
		slog.Warn("Error receiving product information", "status", response.Status, "store_id", c.StoreID, "catalog_id", catalogID)
		return cycletls.Response{}, fmt.Errorf("unexpected status code: %d", response.Status)
	}
	slog.Info("Product info received successfully", "store_id", c.StoreID, "catalog_id", catalogID, "status", response.Status)
	return response, err
}
