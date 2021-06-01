package gout

import (
	"net/http"

	"github.com/guonaihong/gout/dataflow"
)

type Client struct {
	hc *http.Client
	*dataflow.DataFlow
}

func NewWithOpt(opts ...Option) *Client {
	c := &Client{}

	options := options{}

	for _, o := range opts {
		o.apply(&options)
	}

	c.hc = options.hc
	if c.hc == nil {
		c.hc = &http.Client{}
	}

	g := dataflow.New(c.hc)
	c.DataFlow = &g.DataFlow
	return c

}
