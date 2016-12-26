package client

import (
	"container/ring"
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

type Client ring.Ring

func NewClient(cfgs ...ClientCfg) *Client {
	if len(cfgs) == 0 {
		return &Client{}
	}
	r := ring.New(len(cfgs))
	for _, cfg := range cfgs {
		if cfg.proxyURL == "" {
			r.Value = &client{
				&cfg,
				new(http.Client),
				"",
			}
		} else {
			proxyURL, err := url.Parse(cfg.proxyURL)
			if err != nil {
				continue
			}
			r.Value = &client{
				&cfg,
				&http.Client{
					Transport: &http.Transport{
						Proxy: http.ProxyURL(proxyURL),
					},
				},
				"",
			}
		}
		r = r.Next()

	}

	ret := Client(*r)
	return &ret
}
