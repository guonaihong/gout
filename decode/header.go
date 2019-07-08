package decode

import (
	"net/http"
	"net/textproto"
)

type headerDecode struct{}

func (h *headerDecode) Decode(r *http.Request, obj interface{}) error {
	return decode(r.Header, obj, "header")
}

type headerSet map[string][]string

func (h *headerSet) Set(value reflect.Value, sf reflect.StructField, tagValue string) {
	setForm(*h, value, sf, textproto.CanonicalMIMEHeaderKey(tagValue))
}
