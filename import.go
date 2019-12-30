package gout

import (
	"bufio"
	"bytes"
	"github.com/guonaihong/gout/core"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Import struct{}

func NewImport() *Import {
	return &Import{}
}

const rawTextSpace = "\r\n "

func (i *Import) RawText(text interface{}) *Text {
	var read io.Reader
	r := &Text{}

	out := New()
	// TODO 函数
	r.DataFlow.out = out
	r.DataFlow.Req.g = out

	switch body := text.(type) {
	case string:
		body = strings.TrimLeft(body, rawTextSpace)
		read = strings.NewReader(body)
	case []byte:
		body = bytes.TrimLeft(body, rawTextSpace)
		read = bytes.NewReader(body)
	default:
		r.err = core.ErrUnknownType
		return r
	}

	req, err := http.ReadRequest(bufio.NewReader(read))
	if err != nil {
		r.err = err
		return r
	}

	// TODO 探索下能否支持https
	req.URL.Scheme = "http"
	req.URL.Host = req.Host
	req.RequestURI = ""

	if req.GetBody == nil {
		all, err := ioutil.ReadAll(req.Body)
		if err != nil {
			r.err = err
			return r
		}

		req.GetBody = func() (io.ReadCloser, error) {
			return ioutil.NopCloser(bytes.NewReader(all)), nil
		}

		req.Body = ioutil.NopCloser(bytes.NewReader(all))
	}

	r.setRequest(req)

	return r
}
