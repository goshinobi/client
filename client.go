package client

import (
	"fmt"
	"net/http"
	"net/url"

	"h12.me/socks"

	tor "github.com/goshinobi/tor_multi"
)

type ClientCfg struct {
	ProxyType string
	UserAgent string
	UseNUM    int
	MaxTTL    int
	ProxyURL  string
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
	if cfg.ProxyURL == "" {
		return &client{
			&cfg,
			new(http.Client),
			"",
		}, nil
	}
	proxyURL, err := url.Parse(cfg.ProxyURL)
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

func newTorClient(info *tor.ProxyInfo) (*client, error) {
	dialSocksProxy := socks.DialSocksProxy(socks.SOCKS5, fmt.Sprintf("127.0.0.1:%v", info.Conf.SocksPort))
	tr := &http.Transport{Dial: dialSocksProxy}
	httpClient := &http.Client{Transport: tr}

	return &client{
		&ClientCfg{ProxyType: "tor"},
		httpClient,
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

func NewClientTor(n int) *Client {
	ret := Client{
		0,
		make([]*client, 0, n),
	}
	for i := 0; i < n; i++ {
		err := tor.StartProxy()
		if err != nil {
			continue
		}
	}
	proxyList := tor.GetWorkProxyList()
	for _, v := range proxyList {
		torClient, err := newTorClient(v)
		if err != nil {
			continue
		}
		ret.list = append(ret.list, torClient)
	}

	return &ret
}

func (c *Client) Len() int {
	return len(c.list)
}

func (c *Client) next() {
	c.p++
	if c.Len() <= c.p {
		c.p = 0
	}
}

func (cBuffer *Client) Do(req *http.Request) (*http.Response, error) {
	defer func() {
		cBuffer.next()
	}()

	return cBuffer.list[cBuffer.p].c.Do(req)
}

func (cBuffer *Client) Get(u string) (*http.Response, error) {
	defer func() {
		cBuffer.next()
	}()
	return cBuffer.list[cBuffer.p].c.Get(u)
}

func (cBuffer *Client) Add(newClients *Client) {
	cBuffer.list = append(cBuffer.list, newClients.list...)
}
