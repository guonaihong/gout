package gout

import (
	"net/http"
)

type Gout struct {
	*http.Client
	routerGroup
}

func New(c *http.Client) *Gout {
	out := &Gout{Client: c}
	if c == nil {
		out.Client = &http.Client{}
	}

	out.routerGroup.out = out
	return out
}
