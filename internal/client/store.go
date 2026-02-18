package client

import (
	"fmt"
	"log/slog"
	"net/url"

	"github.com/Danny-Dasilva/CycleTLS/cycletls"
)

func (c *ClientWithCookies) GetStore(city string) (cycletls.Response, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return cycletls.Response{}, fmt.Errorf("failed to parse base URL: %w", err)
	}

	u = u.JoinPath("wcs", "resources", "mobihub023", "store", city, "stores")

	slog.Debug("Requesting store", "url", u.String(), "city", city, "store_id", c.StoreID)

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
		slog.Error("Error retrieving store", "error", err, "city", city)
		return cycletls.Response{}, err
	}
	if response.Status != 200 {
		slog.Warn("Error receiving Store", "status", response.Status, "city", city)
		return cycletls.Response{}, fmt.Errorf("unexpected status code: %d", response.Status)
	}
	slog.Info("Store received successfully", "city", city, "status", response.Status)
	return response, err
}

func (c *ClientWithCookies) GetStoreByCoords(longitude, latitude float64) (cycletls.Response, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		slog.Error("Failed to parse base URL", "error", err)
		return cycletls.Response{}, fmt.Errorf("failed to parse base URL: %w", err)
	}

	u = u.JoinPath("wcs", "resources", "mobihub023", "store", "13151", "address", "find")
	q := u.Query()
	q.Set("longitude", fmt.Sprintf("%f", longitude))
	q.Set("latitude", fmt.Sprintf("%f", latitude))
	u.RawQuery = q.Encode()

	slog.Debug("Requesting store by coords", "url", u.String(), "latitude", latitude, "longitude", longitude, "store_id", c.StoreID)
	response, err := c.Do(u.String(), cycletls.Options{
		Ja3:       c.ja3,
		UserAgent: c.UA,
		Proxy:     c.Proxy,

		Headers: map[string]string{
			"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
			"Accept-Language":           "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7",
			"Accept-Encoding":           "gzip, deflate, br",
			"Cache-Control":             "max-age=0",
			"Upgrade-Insecure-Requests": "1",
		},
		Timeout:               30,
		EnableConnectionReuse: true,
	},
		"GET",
	)
	if err != nil {
		slog.Error("Error changing the Store", "error", err, "latitude", latitude, "longitude", longitude)
		return cycletls.Response{}, err
	}
	if response.Status != 200 {
		slog.Warn("Error changing the Store", "status", response.Status, "latitude", latitude, "longitude", longitude)
		return cycletls.Response{}, fmt.Errorf("unexpected status code: %d", response.Status)
	}

	slog.Info("Store by coords received successfully", "status", response.Status, "latitude", latitude, "longitude", longitude)

	return response, err
}
