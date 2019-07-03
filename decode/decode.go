package decode

import (
	"net/http"
)

type Set interface {
	Set(r http.Request, obj interface{}) error
	Name() string
}
