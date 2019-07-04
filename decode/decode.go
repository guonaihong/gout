package decode

import (
	"net/http"
)

type Decoder interface {
	Decode(*http.Request, interface{})
}
