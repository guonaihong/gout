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
	*Req
	reqs []*Req
}

func New(c *http.Client) *Gout {
	out := &Gout{Client: c}
	if c == nil {
		out.Client = &http.Client{}
	}

	return out
}

func (g *Gout) GET(url string) *Gout {
	g.Req = NewReq(GET, url, g)
	return g
}

func (g *Gout) POST(url string) *Gout {
	g.Req = NewReq(POST, url, g)
	return g
}

func (g *Gout) PUT(url string) *Gout {
	g.Req = NewReq(PUT, url, g)
	return g
}

func (g *Gout) DELETE(url string) *Gout {
	g.Req = NewReq(DELETE, url, g)
	return g
}

func (g *Gout) PATCH(url string) *Gout {
	g.Req = NewReq(PATCH, url, g)
	return g
}

func (g *Gout) HEAD(url string) *Gout {
	g.Req = NewReq(HEAD, url, g)
	return g
}

func (g *Gout) OPTIONS(url string) *Gout {
	g.Req = NewReq(OPTIONS, url, g)
	return g
}

func (g *Gout) Do() (err error) {
	if g.Req != nil {
		if err = g.Req.Do(); err != nil {
			return err
		}
	}

	for _, r := range g.reqs {
		if err = r.Do(); err != nil {
			return err
		}
	}

	return nil
}

func (g *Gout) Next() *Gout {
	g.reqs = append(g.reqs, g.Req)
	g.Req = nil
	return g
}
