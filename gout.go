package gout

import (
	"net/http"
)

type Gout struct {
	*http.Client
	DataFlow

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

	out.DataFlow.out = out
	return out
}

// default
func Def() *Gout {
	return New()
}

func GET(url string) *DataFlow {
	return New().GET(url)
}

func POST(url string) *DataFlow {
	return New().POST(url)
}

func PUT(url string) *DataFlow {
	return New().PUT(url)
}

func DELETE(url string) *DataFlow {
	return New().DELETE(url)
}

func PATCH(url string) *DataFlow {
	return New().PATCH(url)
}

func HEAD(url string) *DataFlow {
	return New().HEAD(url)
}

func OPTIONS(url string) *DataFlow {
	return New().OPTIONS(url)
}
