package encode

import (
	"net/http"
	"net/url"
)

type QueryEncode struct {
	values url.Values
	r      *http.Request
}

func NewQueryEncode(req *http.Request) *QueryEncode {
	return &QueryEncode{values: make(url.Values)}
}

func (q *QueryEncode) Add(key, val string) error {
	q.values.Add(key, val)
	return nil
}

func (q *QueryEncode) Name() string {
	return "query"
}

func (q *QueryEncode) End() string {
	return q.values.Encode()
}
