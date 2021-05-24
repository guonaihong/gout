package gout

import (
	"net/http"

	"github.com/guonaihong/gout/dataflow"
)

type Client struct {
	hc http.Client
	*dataflow.DataFlow
}

func NewWithOpt(opts ...Option) *Client {
	c := &Client{}

	g := dataflow.New(&c.hc)

	c.DataFlow = &g.DataFlow

	options := options{
		hc: &c.hc,
	}

	for _, o := range opts {
		o.apply(&options)
	}
	return c

}
