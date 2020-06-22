package dataflow

import "net/http"

type RequestUser interface {
	ModifyRequest(req *http.Request) *http.Request
}
