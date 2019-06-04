package gout

import (
	"net/http"
)

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
)

type Gout struct {
	*http.Client
	method string
	url    string
}

func New(c *http.Client) *Gout {
	out := &Gout{Client: c}
	if c == nil {
		out.c = &http.Client{}
	}

	return out
}

func (g *Gout) GET(url string) *Gout {
	g.method = GET
	g.url = url
	return g
}

func (g *Gout) POST(url string) *Gout {
	g.method = POST
	g.url = url
	return g
}

func (g *Gout) PUT(url string) *Gout {
	g.method = PUT
	g.url = url
	return g
}

func (g *Gout) DELETE(url string) *Gout {
	g.method = DELETE
	g.url = url
	return g
}

func (g *Gout) PATCH(url string) *Gout {
	g.method = PATCH
	g.url = url
	return g
}

func (g *Gout) HEAD(url string) *Gout {
	g.method = HEAD
	g.url = url
	return g
}

func (g *Gout) OPTIONS(url string) *Gout {
	g.method = OPTIONS
	g.url = url
	return g
}

func (g *Gout) Do() *Gout {
}
