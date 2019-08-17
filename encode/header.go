package encode

import (
	"net/http"
	"reflect"
)

var _ Adder = (*HeaderEncode)(nil)

type HeaderEncode struct {
	r *http.Request
}

func NewHeaderEncode(req *http.Request) *HeaderEncode {
	return &HeaderEncode{r: req}
}

func (h *HeaderEncode) Add(key string, v reflect.Value, sf reflect.StructField) error {
	h.r.Header.Add(key, valToStr(v, sf))
	return nil
}

func (h *HeaderEncode) Name() string {
	return "header"
}
