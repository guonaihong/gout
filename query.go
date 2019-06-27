package gout

import (
	"net/http"
	"net/url"
)

type QueryEncode struct {
	URL *url.URL
}

func NewQueryEncode(req *http.Request) *QueryEncode {
	return &QueryEncode{URL: req.URL}
}

func (q *QueryEncode) Add(key, val string) error {
	q.URL.Add(key, val)
}

func (q *QueryEncode) Name() string {
	return "query"
}
