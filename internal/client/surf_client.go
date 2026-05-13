package client

import (
	"TestCenozavr/internal/config"
	"log/slog"
	"time"

	"github.com/enetx/g"
	"github.com/enetx/surf"
)

const baseURL = "https://www.okeydostavka.ru"

type Client struct {
	client  *surf.Client
	StoreID uint64
}

func NewClient(cfg *config.Config) *Client {
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

	return &Client{
		client:  c,
		StoreID: cfg.StoreID,
	}
}

func NewClientWithCookies(cfg *config.Config) (*Client, error) {
	client := NewClient(cfg)
	if err := client.GetCookie(); err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) Close() {
	c.client.Close()
	slog.Info("Client closed")
}

func (c *Client) SetStoreID(storeID uint64) {
	c.StoreID = storeID
	slog.Debug("Store ID updated", "new", storeID)
}

func (c *Client) GET(url string) *surf.Request {
	return c.client.Get(g.String(url))
}

func (c *Client) POST(url string) *surf.Request {
	return c.client.Post(g.String(url))
}
