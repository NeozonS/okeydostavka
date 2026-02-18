package client

import (
	"TestCenozavr/internal/config"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"

	"github.com/Danny-Dasilva/CycleTLS/cycletls"
)

const baseURL = "https://www.okeydostavka.ru"

type ClientWithCookies struct {
	client  cycletls.CycleTLS
	jar     http.CookieJar
	mu      sync.RWMutex
	ja3     string
	UA      string
	Proxy   string
	StoreID uint64
}

func NewClientWithCookies(cfg *config.Config) *ClientWithCookies {
	jar, _ := cookiejar.New(nil)
	client := cycletls.Init()

	return &ClientWithCookies{
		client:  client,
		Proxy:   cfg.Proxy,
		StoreID: cfg.StoreID,
		jar:     jar,
		mu:      sync.RWMutex{},
		ja3:     "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-34-18-51-43-13-45-28-27-65037,4588-29-23-24-25-256-257,0",
		UA:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:147.0) Gecko/20100101 Firefox/147.0",
	}
}

func (c *ClientWithCookies) Do(targetURL string, options cycletls.Options, method string) (cycletls.Response, error) {

	u, err := url.Parse(targetURL)
	if err != nil {
		slog.Error("Failed to parse target URL", "error", err, "url", targetURL)
		return cycletls.Response{}, fmt.Errorf("failed to parse URL: %w", err)
	}

	c.mu.Lock()
	existingCookies := c.jar.Cookies(u)
	c.mu.Unlock()
	if len(existingCookies) > 0 {
		options.Cookies = formatCookies(existingCookies)
		slog.Debug("Using existing cookies", "count", len(existingCookies), "url", u.Host)
	}

	response, err := c.client.Do(targetURL, options, method)
	if err != nil {
		slog.Error("Request failed", "error", err, "url", targetURL, "method", method)
		return response, err
	}

	if len(response.Cookies) > 0 {
		c.mu.Lock()
		c.jar.SetCookies(u, response.Cookies)
		c.mu.Unlock()
		slog.Debug("Cookies saved", "count", len(response.Cookies), "url", u.Host)
	}

	return response, nil
}

func (c *ClientWithCookies) ClearCookies() {
	c.mu.Lock()
	defer c.mu.Unlock()

	newJar, err := cookiejar.New(nil)
	if err != nil {
		slog.Error("Failed to clear cookies", "error", err)
		return
	}

	c.jar = newJar
	slog.Info("Cookies cleared")
}

func (c *ClientWithCookies) Close() {
	c.client.Close()
	slog.Info("Client closed")
}

func (c *ClientWithCookies) SetStoreID(storeID uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.StoreID = storeID
	slog.Debug("Store ID updated", "new", storeID)
}

func formatCookies(cookies []*http.Cookie) []cycletls.Cookie {
	out := make([]cycletls.Cookie, 0, len(cookies))

	for _, c := range cookies {
		out = append(out, cycletls.Cookie{
			Name:  c.Name,
			Value: c.Value,
		})
	}
	return out
}
