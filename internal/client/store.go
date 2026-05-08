package client

import (
	"fmt"
	"log/slog"
	"net/url"
)

func (c *ClientWithCookies) GetStore(city string) ([]byte, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	u = u.JoinPath("wcs", "resources", "mobihub023", "store", city, "stores")

	slog.Debug("Requesting store", "url", u.String(), "city", city, "store_id", c.StoreID)

	resp := c.GET(u.String()).
		AddHeaders("Accept", "application/json, text/plain, */*").
		AddHeaders("Content-Type", "application/json").
		AddHeaders("Origin", baseURL).
		AddHeaders("X-Requested-With", "XMLHttpRequest").
		Do()

	if resp.IsErr() {
		slog.Error("Error retrieving store", "error", resp.Err(), "city", city)
		return nil, resp.Err()
	}

	r := resp.Ok()
	if r.StatusCode != 200 {
		slog.Warn("Error receiving Store", "status", r.StatusCode, "city", city)
		return nil, fmt.Errorf("unexpected status code: %d", r.StatusCode)
	}

	body := r.Body.Bytes()
	if body.IsErr() {
		return nil, body.Err()
	}

	slog.Info("Store received successfully", "city", city, "status", r.StatusCode)
	return body.Ok(), nil
}

func (c *ClientWithCookies) GetStoreByCoords(longitude, latitude float64) ([]byte, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		slog.Error("Failed to parse base URL", "error", err)
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	u = u.JoinPath("wcs", "resources", "mobihub023", "store", "13151", "address", "find")
	q := u.Query()
	q.Set("longitude", fmt.Sprintf("%f", longitude))
	q.Set("latitude", fmt.Sprintf("%f", latitude))
	u.RawQuery = q.Encode()

	slog.Debug("Requesting store by coords", "url", u.String(), "latitude", latitude, "longitude", longitude, "store_id", c.StoreID)

	resp := c.GET(u.String()).
		SetHeaders("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7").
		AddHeaders("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8").
		AddHeaders("Cache-Control", "max-age=0").
		AddHeaders("Upgrade-Insecure-Requests", "1").
		Do()

	if resp.IsErr() {
		slog.Error("Error changing the Store", "error", resp.Err(), "latitude", latitude, "longitude", longitude)
		return nil, resp.Err()
	}

	r := resp.Ok()
	if r.StatusCode != 200 {
		slog.Warn("Error changing the Store", "status", r.StatusCode, "latitude", latitude, "longitude", longitude)
		return nil, fmt.Errorf("unexpected status code: %d", r.StatusCode)
	}

	body := r.Body.Bytes()
	if body.IsErr() {
		return nil, body.Err()
	}

	slog.Info("Store by coords received successfully", "status", r.StatusCode, "latitude", latitude, "longitude", longitude)
	return body.Ok(), nil
}
