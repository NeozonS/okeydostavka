package client

import (
	"fmt"
	"log/slog"
)

func (c *ClientWithCookies) GetCookie() error {
	resp := c.GET(baseURL).
		SetHeaders("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7").
		AddHeaders("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8").
		AddHeaders("Referer", "https://www.okeydostavka.ru/spb").
		AddHeaders("Cache-Control", "max-age=0").
		AddHeaders("Upgrade-Insecure-Requests", "1").
		Do()

	if resp.IsErr() {
		slog.Error("Cookie receipt error", "error", resp.Err())
		return fmt.Errorf("cookie request failed: %w", resp.Err())
	}

	r := resp.Ok()
	if r.StatusCode != 200 {
		slog.Warn("Cookie receipt error", "status", r.StatusCode)
		r.Debug().Request().Response(true).Print()
		return fmt.Errorf("unexpected status code: %d", r.StatusCode)
	}

	slog.Info("Cookies received successfully", "status", r.StatusCode)
	return nil
}
