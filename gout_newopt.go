package gout

import (
	"net/http"

	"github.com/guonaihong/gout/dataflow"
)

type Client struct {
	options
}

func NewWithOpt(opts ...Option) *Client {
	c := &Client{}
	c.hc = &http.Client{}

	for _, o := range opts {
		o.apply(&c.options)
	}

	return c

}

// GET send HTTP GET method
func (c *Client) GET(url string) *dataflow.DataFlow {
	return dataflow.New(c.hc).GET(url)
}

// POST send HTTP POST method
func (c *Client) POST(url string) *dataflow.DataFlow {
	return dataflow.New(c.hc).POST(url)
}

// PUT send HTTP PUT method
func (c *Client) PUT(url string) *dataflow.DataFlow {
	return dataflow.New(c.hc).PUT(url)
}

// DELETE send HTTP DELETE method
func (c *Client) DELETE(url string) *dataflow.DataFlow {
	return dataflow.New(c.hc).DELETE(url)
}

// PATCH send HTTP PATCH method
func (c *Client) PATCH(url string) *dataflow.DataFlow {
	return dataflow.New(c.hc).PATCH(url)
}

// HEAD send HTTP HEAD method
func (c *Client) HEAD(url string) *dataflow.DataFlow {
	return dataflow.New(c.hc).HEAD(url)
}

// OPTIONS send HTTP OPTIONS method
func (c *Client) OPTIONS(url string) *dataflow.DataFlow {
	return dataflow.New(c.hc).OPTIONS(url)
}
