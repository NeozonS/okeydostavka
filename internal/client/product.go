package client

import (
	"fmt"
	"log/slog"
	"net/url"
)

func (c *ClientWithCookies) GetProductInfo(catalogID string, pageNumber, pageSize int) ([]byte, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		slog.Error("Failed to parse base URL", "error", err)
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	u = u.JoinPath("wcs", "resources", "mobihub023", "store", fmt.Sprintf("%d", c.StoreID), "catalog", "combo", catalogID)
	q := u.Query()
	q.Set("pageSize", fmt.Sprintf("%d", pageSize))
	q.Set("pageNumber", fmt.Sprintf("%d", pageNumber))
	u.RawQuery = q.Encode()

	resp := c.POST(u.String()).
		SetHeaders("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7").
		AddHeaders("Referer", baseURL+"/").
		AddHeaders("Accept", "*/*").
		AddHeaders("Content-Language", "ru-RU").
		AddHeaders("Content-Type", "application/x-www-form-urlencoded").
		AddHeaders("Priority", "u=4").
		AddHeaders("X-Requested-With", "XMLHttpRequest").
		AddHeaders("Sec-Fetch-Dest", "empty").
		AddHeaders("Sec-Fetch-Mode", "cors").
		AddHeaders("Sec-Fetch-Site", "same-origin").
		Do()

	if resp.IsErr() {
		slog.Error("Price receipt request error", "error", resp.Err(), "store_id", c.StoreID, "catalog_id", catalogID)
		return nil, resp.Err()
	}

	r := resp.Ok()
	if r.StatusCode == 403 {
		return nil, AntiBotError
	}
	if r.StatusCode != 200 {
		slog.Warn("Error receiving product information", "status", r.StatusCode, "store_id", c.StoreID, "catalog_id", catalogID)
		return nil, fmt.Errorf("unexpected status code: %d", r.StatusCode)
	}

	body := r.Body.Bytes()
	if body.IsErr() {
		return nil, body.Err()
	}

	slog.Info("Product info received successfully", "store_id", c.StoreID, "catalog_id", catalogID, "status", r.StatusCode)
	return body.Ok(), nil
}
