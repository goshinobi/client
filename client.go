package client

import (
	"net/http"
	"time"
)

type client struct {
	c          *http.Client
	lifespan   int
	usedNum    int
	ip         string
	useragent  string
	errorcount int
}

func (c *client) do(req *http.Request) (resp *http.Response, err error) {
	c.lifespan++
	return c.c.Do(req)
}

func (c *client) isDie() bool {
	return (c.errorcount > 5) || ((c.usedNum > c.lifespan) && (c.usedNum >= 0))
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

func (c *Client) Add(h *http.Client, lifespan int) {
	c.clients = append(c.clients, &client{
		c:        h,
		lifespan: lifespan,
	})
}

func New() *Client {
	return &Client{}
}
