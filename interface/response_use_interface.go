package dataflow

import (
	"net/http"
)

// ResponseMiddler 响应拦截器
type ResponseMiddler interface {
	ModifyResponse(response *http.Response) error
}
