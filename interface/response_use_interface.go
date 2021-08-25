package dataflow

import (
	"net/http"
)

// 响应拦截器
type ResponseMiddler interface {
	ModifyResponse(response *http.Response) error
}
