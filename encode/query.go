package encode

import (
	"net/http"
	"net/url"
	"reflect"
)

var _ Adder = (*QueryEncode)(nil)

type QueryEncode struct {
	values url.Values
	r      *http.Request
}

func NewQueryEncode(req *http.Request) *QueryEncode {
	return &QueryEncode{values: make(url.Values)}
}

func (q *QueryEncode) Add(key string, v reflect.Value, sf reflect.StructField) error {
	q.values.Add(key, valToStr(v, sf))
	return nil
}

func (q *QueryEncode) Name() string {
	return "query"
}

func (q *QueryEncode) End() string {
	return q.values.Encode()
}
