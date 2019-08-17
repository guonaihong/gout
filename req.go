package gout

import (
	"bytes"
	"fmt"
	"github.com/guonaihong/gout/decode"
	"github.com/guonaihong/gout/encode"
	"net/http"
	"reflect"
	"strings"
)

type Req struct {
	method string
	url    string

	formEncode interface{}

	// http body
	bodyEncoder Encoder
	bodyDecoder Decoder

	// http header
	headerEncode interface{}
	headerDecode interface{}

	// query
	queryEncode interface{}

	httpCode *int
	g        *Gout
}

func (r *Req) Reset() {
	r.formEncode = nil
	r.bodyEncoder = nil
	r.bodyDecoder = nil
	r.httpCode = nil
	r.headerDecode = nil
	r.headerEncode = nil
	r.queryEncode = nil
}

func isString(x interface{}) (string, bool) {
	p := reflect.ValueOf(x)

	for p.Kind() == reflect.Ptr {
		p = p.Elem()
	}

	if p.Kind() == reflect.String {
		s := p.Interface().(string)
		if strings.HasPrefix(s, "?") {
			s = s[1:]
		}
		return s, true
	}
	return "", false
}

func (r *Req) Do() (err error) {
	b := &bytes.Buffer{}

	defer r.Reset()

	// set http body
	if r.bodyEncoder != nil {
		if err := r.bodyEncoder.Encode(b); err != nil {
			return err
		}
	}

	// set query header
	if r.queryEncode != nil {
		var query string
		if q, ok := isString(r.queryEncode); ok {
			query = q
		} else {
			q := encode.NewQueryEncode(nil)
			if err = encode.Encode(r.queryEncode, q); err != nil {
				return err
			}

			query = q.End()
		}

		if len(query) > 0 {
			r.url += "?" + query
		}
	}

	var f *encode.FormEncode

	if r.formEncode != nil {
		f = encode.NewFormEncode(b)
		if err := encode.Encode(r.formEncode, f); err != nil {
			return err
		}

		f.End()
	}

	req, err := http.NewRequest(r.method, r.url, b)
	if err != nil {
		return err
	}

	if r.formEncode != nil {
		req.Header.Add("Content-Type", f.FormDataContentType())
	}

	// set http header
	if r.headerEncode != nil {
		err = encode.Encode(r.headerEncode, encode.NewHeaderEncode(req))
		if err != nil {
			return err
		}
	}

	resp, err := r.g.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if r.headerDecode != nil {
		err = decode.Header.Decode(resp, r.headerDecode)
		if err != nil {
			return err
		}
	}

	if r.bodyDecoder != nil {
		if err := r.bodyDecoder.Decode(resp.Body); err != nil {
			return err
		}
	}

	if r.httpCode != nil {
		*r.httpCode = resp.StatusCode
	}

	return nil
}

func modifyUrl(url string) string {
	if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://") {
		return url
	}

	if strings.HasPrefix(url, ":") {
		return fmt.Sprintf("http://127.0.0.1%s", url)
	}

	if strings.HasPrefix(url, "/") {
		return fmt.Sprintf("http://127.0.0.1%s", url)
	}

	return fmt.Sprintf("http://%s", url)
}

func NewReq(method string, url string, g *Gout) *Req {
	return &Req{method: method, url: modifyUrl(url), g: g}
}
