package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
)

type gzipCompress struct{}

func (g *gzipCompress) ModifyRequest(req *http.Request) (*http.Request, error) {
	buf := &bytes.Buffer{}

	w := gzip.NewWriter(buf)
	body, err := req.GetBody()
	if err != nil {
		return nil, err
	}
	w.Close()
	io.Copy(w, body)

	req.Body = ioutil.NopCloser(buf)
	return req, nil
}
