package encode

import (
	"net/http"
	"net/url"
	"reflect"
)

var _ Adder = (*QueryEncode)(nil)

// QueryEncode URL query encoder structure
type QueryEncode struct {
	values url.Values
	r      *http.Request
}

// NewQueryEncode create a new URL query  encoder
func NewQueryEncode(req *http.Request) *QueryEncode {
	return &QueryEncode{values: make(url.Values)}
}

// Add Encoder core function, used to set each key / value into the http URL query
func (q *QueryEncode) Add(key string, v reflect.Value, sf reflect.StructField) error {
	q.values.Add(key, valToStr(v, sf))
	return nil
}

// End URL query structured data into strings
func (q *QueryEncode) End() string {
	return q.values.Encode()
}

// Name URL query Encoder name
func (q *QueryEncode) Name() string {
	return "query"
}
