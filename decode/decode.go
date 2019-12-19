package decode

import (
	"net/http"
)

// Decoder Decoder interface
type Decoder interface {
	Decode(*http.Request, interface{})
}

var (
	// Header is the http header decoder
	Header = headerDecode{}
)
