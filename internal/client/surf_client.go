package client

import (
	"TestCenozavr/internal/config"
	"log/slog"
	"time"

	"github.com/enetx/g"
	"github.com/enetx/surf"
)

const baseURL = "https://www.okeydostavka.ru"

type ClientWithCookies struct {
	client  *surf.Client
	StoreID uint64
}

func NewClientWithCookies(cfg *config.Config) *ClientWithCookies {
	c := surf.NewClient().
		Builder().
		Impersonate().Firefox().
		Session().
		Retry(3, 2*time.Second).
		Proxy(g.String(cfg.Proxy)).
		Timeout(30 * time.Second).
		MaxRedirects(10).
		ForceHTTP1().
		Build().
		Unwrap()

	return &ClientWithCookies{
		client:  c,
		StoreID: cfg.StoreID,
	}
}

func (c *ClientWithCookies) Close() {
	c.client.Close()
	slog.Info("Client closed")
}

func (c *ClientWithCookies) SetStoreID(storeID uint64) {
	c.StoreID = storeID
	slog.Debug("Store ID updated", "new", storeID)
}

func (c *ClientWithCookies) GET(url string) *surf.Request {
	return c.client.Get(g.String(url))
}

func (c *ClientWithCookies) POST(url string) *surf.Request {
	return c.client.Post(g.String(url))
}
