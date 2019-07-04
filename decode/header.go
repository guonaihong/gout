package decode

import (
	"net/http"
)

type headerDecode struct{}

func (h *headerDecode) Decode(r *http.Request, obj interface{}) error {
	return decode(r.Header, obj)
}

func decodeHeader(header map[string][]string, obj interface{}) error {
	return decode(header, obj)
}
