package client

import (
	"net/http"
	"net/url"
)

type ClientCfg struct {
	userAgent string
	useNUM    int
	maxTTL    int
	proxyURL  string
}

type client struct {
	cfg *ClientCfg
	c   *http.Client
	ip  string
}

type Client struct {
	p    int
	list []*client
}

func newClient(cfg ClientCfg) (*client, error) {
	if cfg.proxyURL == "" {
		return &client{
			&cfg,
			new(http.Client),
			"",
		}, nil
	}
	proxyURL, err := url.Parse(cfg.proxyURL)
	if err != nil {
		return nil, err
	}
	return &client{
		&cfg,
		&http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		},
		"",
	}, nil
}

func NewClient(cfgs ...ClientCfg) *Client {
	ret := Client{
		0,
		make([]*client, 0, len(cfgs)),
	}
	for _, cfg := range cfgs {
		c, err := newClient(cfg)
		if err != nil {
			continue
		}
		ret.list = append(ret.list, c)
	}
	return &ret
}

func (c *Client) Len() int {
	return len(c.list)
}
