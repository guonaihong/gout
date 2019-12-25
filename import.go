package gout

import (
	"bufio"
	"net/http"
	"strings"
)

type Import struct{}

func NewImport() *Import {
	return &Import{}
}

func (i *Import) RawText(s string) *Text {
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(s)))
	if err != nil {
		r := &Text{}
		r.err = err
		return r
	}

	r := &Text{}
	r.setRequest(req)
	out := New()

	// todo 函数
	r.DataFlow.out = out
	r.DataFlow.Req.g = out
	return r
}
