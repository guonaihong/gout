package gout

import (
	"net/http"
	"net/url"
)

type QueryEncode struct {
	r    *http.Request
	vals url.Values
}

func NewQueryEncode(req *http.Request) *QueryEncode {
	return &QueryEncode{vals: make(url.Values)}
}

func (q *QueryEncode) Add(key, val string) error {
	q.vals.Add(key, val)
	return nil
}

func (q *QueryEncode) Name() string {
	return "query"
}

func (q *QueryEncode) End() {
	_ = q.vals.Encode()
}
