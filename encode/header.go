package encode

import (
	"net/http"
)

type HeaderEncode struct {
	header http.Header
}

func NewHeaderEnocde(req *http.Request) *HeaderEncode {
	return &HeaderEnocde{obj: obj, header: req.Header}
}

func (h *HeaderEncode) Add(key, val string) error {
	h.Add(key, val)
}

func (h *HeaderEncode) Name() string {
	return "header"
}
