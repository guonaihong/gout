package gout

import (
	"net/http"
)

type Gout struct {
	*http.Client
	routerGroup

	opt DebugOption
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

func GET(url string) *routerGroup {
	return New(nil).GET(url)
}

func POST(url string) *routerGroup {
	return New(nil).POST(url)
}

func PUT(url string) *routerGroup {
	return New(nil).PUT(url)
}

func DELETE(url string) *routerGroup {
	return New(nil).DELETE(url)
}

func PATCH(url string) *routerGroup {
	return New(nil).PATCH(url)
}

func HEAD(url string) *routerGroup {
	return New(nil).HEAD(url)
}

func OPTIONS(url string) *routerGroup {
	return New(nil).OPTIONS(url)
}
