package gout

import (
	"context"
	"fmt"
	"github.com/guonaihong/gout/decode"
	"github.com/guonaihong/gout/encode"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	Get     = "GET"
	Post    = "POST"
	Put     = "PUT"
	Delete  = "DELETE"
	Patch   = "PATCH"
	Head    = "HEAD"
	Options = "OPTIONS"
)

type DataFlow struct {
	basePath string
	Req
	out *Gout
}

func (df *DataFlow) GET(url string) *DataFlow {
	df.Req = reqDef(Get, joinPaths(df.basePath, url), df.out)
	return df
}

func (df *DataFlow) POST(url string) *DataFlow {
	df.Req = reqDef(Post, joinPaths(df.basePath, url), df.out)
	return df
}

func (df *DataFlow) PUT(url string) *DataFlow {
	df.Req = reqDef(Put, joinPaths(df.basePath, url), df.out)
	return df
}

func (df *DataFlow) DELETE(url string) *DataFlow {
	df.Req = reqDef(Delete, joinPaths(df.basePath, url), df.out)
	return df
}

func (df *DataFlow) PATCH(url string) *DataFlow {
	df.Req = reqDef(Patch, joinPaths(df.basePath, url), df.out)
	return df
}

func (df *DataFlow) HEAD(url string) *DataFlow {
	df.Req = reqDef(Head, joinPaths(df.basePath, url), df.out)
	return df
}

func (df *DataFlow) OPTIONS(url string) *DataFlow {
	df.Req = reqDef(Options, joinPaths(df.basePath, url), df.out)
	return df
}

func (df *DataFlow) SetBody(obj interface{}) *DataFlow {
	df.Req.bodyEncoder = encode.NewBodyEncode(obj)
	return df
}

func (df *DataFlow) SetForm(obj interface{}) *DataFlow {
	df.Req.formEncode = obj
	return df
}

func (df *DataFlow) SetWWWForm(obj interface{}) *DataFlow {
	df.Req.bodyEncoder = encode.NewWWWFormEncode(obj)
	return df
}

func (df *DataFlow) SetQuery(obj interface{}) *DataFlow {
	df.Req.queryEncode = obj
	return df
}

func (df *DataFlow) SetHeader(obj interface{}) *DataFlow {
	df.Req.headerEncode = obj
	return df
}

func (df *DataFlow) SetJSON(obj interface{}) *DataFlow {
	df.out.opt.ReqBodyType = "json"
	df.Req.bodyEncoder = encode.NewJSONEncode(obj)
	return df
}

func (df *DataFlow) SetXML(obj interface{}) *DataFlow {
	df.out.opt.ReqBodyType = "xml"
	df.Req.bodyEncoder = encode.NewXMLEncode(obj)
	return df
}

func (df *DataFlow) SetYAML(obj interface{}) *DataFlow {
	df.out.opt.ReqBodyType = "yaml"
	df.Req.bodyEncoder = encode.NewYAMLEncode(obj)
	return df
}

// 此函数会修改Transport
func (df *DataFlow) UnixSocket(path string) *DataFlow {
	if df.out.Client.Transport == nil {
		df.out.Client.Transport = &http.Transport{}
	}

	transport, ok := df.out.Client.Transport.(*http.Transport)
	if !ok {
		df.Req.err = fmt.Errorf("UnixSocket:not found http.transport:%T", df.out.Client.Transport)
		return df
	}

	transport.Dial = func(proto, addr string) (conn net.Conn, err error) {
		return net.Dial("unix", path)
	}

	return df
}

// 此函数会修改Transport
func (df *DataFlow) SetProxy(proxyURL string) *DataFlow {
	proxy, err := url.Parse(modifyURL(proxyURL))
	if err != nil {
		df.Req.err = err
		return df
	}

	if df.out.Client.Transport == nil {
		df.out.Client.Transport = &http.Transport{}
	}

	transport, ok := df.out.Client.Transport.(*http.Transport)
	if !ok {
		df.Req.err = fmt.Errorf("SetProxy:not found http.transport:%T", df.out.Client.Transport)
		return df
	}

	transport.Proxy = http.ProxyURL(proxy)

	return df
}

func (df *DataFlow) SetCookies(c ...*http.Cookie) *DataFlow {
	df.Req.cookies = append(df.Req.cookies, c...)
	return df
}

func (df *DataFlow) BindBody(obj interface{}) *DataFlow {
	df.Req.bodyDecoder = decode.NewBodyDecode(obj)
	return df
}

func (df *DataFlow) BindHeader(obj interface{}) *DataFlow {
	df.Req.headerDecode = obj
	return df
}

func (df *DataFlow) BindJSON(obj interface{}) *DataFlow {
	df.out.opt.RspBodyType = "json"
	df.Req.bodyDecoder = decode.NewJSONDecode(obj)
	return df
}

func (df *DataFlow) BindXML(obj interface{}) *DataFlow {
	df.out.opt.RspBodyType = "xml"
	df.Req.bodyDecoder = decode.NewXMLDecode(obj)
	return df
}

func (df *DataFlow) BindYAML(obj interface{}) *DataFlow {
	df.out.opt.RspBodyType = "yaml"
	df.Req.bodyDecoder = decode.NewYAMLDecode(obj)
	return df
}

func (df *DataFlow) Code(httpCode *int) *DataFlow {
	df.Req.httpCode = httpCode
	return df
}

func (df *DataFlow) Callback(cb func(*Context) error) *DataFlow {
	df.Req.callback = cb
	return df
}

func (df *DataFlow) SetTimeout(d time.Duration) *DataFlow {
	df.Req.index++
	df.Req.timeoutIndex = df.Req.index
	df.Req.timeout = d
	return df
}

func (df *DataFlow) WithContext(c context.Context) *DataFlow {
	df.Req.index++
	df.Req.ctxIndex = df.Req.index
	df.Req.c = c
	return df
}

func (df *DataFlow) Debug(d ...interface{}) *DataFlow {
	for _, v := range d {
		switch opt := v.(type) {
		case bool:
			defaultDebug(&df.out.opt)
		case DebugOpt:
			opt.Apply(&df.out.opt)
		}
	}

	return df
}

func (df *DataFlow) Do() (err error) {
	return df.Req.Do()
}

func (df *DataFlow) Filter() *Filter {
	return &Filter{df: df}
}
