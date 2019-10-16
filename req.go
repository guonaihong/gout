package gout

import (
	"bytes"
	"context"
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

	callback func(*Context) error

	//cookie
	cookies []*http.Cookie

	c   context.Context
	err error
}

// req 结构布局说明，以decode为例
// body 可以支持text, json, yaml, xml，所以定义成接口形式
// headerDecode只有一个可能，就定义为具体类型。这里他们的decode实现也不一样
// 有没有必要，归一化成一种???

func (r *Req) Reset() {
	r.err = nil
	r.cookies = nil
	r.formEncode = nil
	r.bodyEncoder = nil
	r.bodyDecoder = nil
	r.httpCode = nil
	r.headerDecode = nil
	r.headerEncode = nil
	r.queryEncode = nil
	r.c = nil
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
	if r.err != nil {
		return r.err
	}

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

	if r.c != nil {
		req = req.WithContext(r.c)
	}

	for _, c := range r.cookies {
		req.AddCookie(c)
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

	fmt.Printf("r.g.opt = %p\n", &r.g.opt)
	if r.g.opt.Debug {
		if err := r.g.opt.resetBodyAndPrint(req, resp); err != nil {
			return err
		}
	}

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

	if r.callback != nil {
		c := Context{Code: resp.StatusCode, Resp: resp}
		if err := r.callback(&c); err != nil {
			return err
		}
	}

	return nil
}

func modifyURL(url string) string {
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
	return &Req{method: method, url: modifyURL(url), g: g}
}
