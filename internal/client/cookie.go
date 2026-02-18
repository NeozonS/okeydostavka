package client

import (
	"fmt"
	"log/slog"

	"github.com/Danny-Dasilva/CycleTLS/cycletls"
)

func (c *ClientWithCookies) GetCookie() error {
	response, err := c.Do(baseURL, cycletls.Options{
		Ja3:       c.ja3,
		UserAgent: c.UA,
		Proxy:     c.Proxy,

		Headers: map[string]string{
			"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
			"Accept-Language":           "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7",
			"Accept-Encoding":           "gzip, deflate, br",
			"Connection":                "keep-alive",
			"Referer":                   "https://www.okeydostavka.ru/spb",
			"Cache-Control":             "max-age=0",
			"Upgrade-Insecure-Requests": "1",
		},
		Timeout:               30,
		EnableConnectionReuse: true,
	}, "GET")

	if err != nil {
		slog.Error("Cookie receipt error", "error", err)
		return fmt.Errorf("cookie request failed: %w", err)
	}
	if response.Status != 200 {
		slog.Warn("Cookie receipt error", "status", response.Status)
		return fmt.Errorf("unexpected status code: %d", response.Status)
	}
	slog.Info("Cookies received successfully", "status", response.Status)

	return nil
}
