package encode

import (
	"net/http"
)

type HeaderEncode struct {
	r *http.Request
}

func NewHeaderEnocde(req *http.Request) *HeaderEncode {
	return &HeaderEncode{r: req}
}

func (h *HeaderEncode) Add(key, val string) error {
	h.r.Header.Add(key, val)
	return nil
}

func (h *HeaderEncode) Name() string {
	return "header"
}
