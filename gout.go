package gout

import (
	"net/http"
)

type Gout struct {
	*http.Client
	routerGroup

	opt DebugOption
}

var (
	DefaultClient      = http.Client{}
	DefaultBenchClient = http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10000,
		},
	}
)

func New(c ...*http.Client) *Gout {
	out := &Gout{}
	if len(c) == 0 || c[0] == nil {
		out.Client = &DefaultClient
	} else {
		out.Client = c[0]
	}

	out.routerGroup.out = out
	return out
}

// default
func Def() *Gout {
	return New()
}

func GET(url string) *routerGroup {
	return New().GET(url)
}

func POST(url string) *routerGroup {
	return New().POST(url)
}

func PUT(url string) *routerGroup {
	return New().PUT(url)
}

func DELETE(url string) *routerGroup {
	return New().DELETE(url)
}

func PATCH(url string) *routerGroup {
	return New().PATCH(url)
}

func HEAD(url string) *routerGroup {
	return New().HEAD(url)
}

func OPTIONS(url string) *routerGroup {
	return New().OPTIONS(url)
}
