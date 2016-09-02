package client

import (
	"errors"
	"net/http"
	"time"

	"github.com/goshinobi/client/ipchecker"
)

type client struct {
	c          *http.Client
	conf       *Config
	ip         string
	errorcount int
}

func (c *client) do(req *http.Request) (resp *http.Response, err error) {
	if c.conf != nil {
		c.conf.lifespan++
	}
	return c.c.Do(req)
}

func (c *client) isDie() bool {
	if c.conf == nil {
		return false
	}
	return (c.errorcount > 5) || ((c.conf.usedNum > c.conf.lifespan) && (c.conf.usedNum >= 0))
}

type Config struct {
	lifespan  int
	usedNum   int
	useragent string
}

type Client struct {
	clients         []*client
	maxResponseTime time.Time
	checkOwnIP      bool
	cn              int
}

func (c *Client) deleteClient(n int) {
	c.clients = append(c.clients[:n], c.clients[n+1:]...)
}

func (c *Client) Do(req *http.Request) (resp *http.Response, err error) {
	if len(c.clients) == 0 {
		return nil, errors.New("no set clients")
	}
	resp, err = c.clients[c.cn].do(req)
	if err != nil {
		c.clients[c.cn].errorcount++
	}
	defer func() {
		if c.clients[c.cn].isDie() {
			c.deleteClient(c.cn)
		}
		c.cn++
	}()
	return resp, err
}

func (c *Client) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Add(h *http.Client, conf *Config) {
	var ip string
	if c.checkOwnIP {
		ip = ipchecker.Check(h)
	}
	c.clients = append(c.clients, &client{
		c:    h,
		conf: conf,
		ip:   ip,
	})
}

func (c *Client) GetIPs() []string {
	var result []string
	if !c.checkOwnIP {
		return nil
	}
	for _, client := range c.clients {
		result = append(result, client.ip)
	}
	return result
}

func New(checkOwnIP bool) *Client {
	return &Client{checkOwnIP: checkOwnIP}
}
