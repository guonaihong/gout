package gout

import (
	"net/http"
)

type Gout struct {
	*http.Client
	routerGroup
}

var DefaultClient = http.Client{}

func New(c *http.Client) *Gout {
	out := &Gout{Client: c}
	if c == nil {
		out.Client = &DefaultClient
	}

	out.routerGroup.out = out
	return out
}

// default
func Def() *Gout {
	return New(nil)
}
