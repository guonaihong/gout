package gout

import (
	"context"
	"fmt"
	"github.com/guonaihong/gout/decode"
	"github.com/guonaihong/gout/encode"
	"net"
	"net/http"
	"net/url"
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

type routerGroup struct {
	basePath string
	*Req
	reqs []*Req
	out  *Gout
}

func (g *routerGroup) Group(relativePath string) *routerGroup {
	return &routerGroup{
		basePath: joinPaths(g.basePath, relativePath),
		out:      g.out,
	}
}

func (g *routerGroup) GET(url string) *routerGroup {
	g.Req = NewReq(Get, joinPaths(g.basePath, url), g.out)
	return g
}

func (g *routerGroup) POST(url string) *routerGroup {
	g.Req = NewReq(Post, joinPaths(g.basePath, url), g.out)
	return g
}

func (g *routerGroup) PUT(url string) *routerGroup {
	g.Req = NewReq(Put, joinPaths(g.basePath, url), g.out)
	return g
}

func (g *routerGroup) DELETE(url string) *routerGroup {
	g.Req = NewReq(Delete, joinPaths(g.basePath, url), g.out)
	return g
}

func (g *routerGroup) PATCH(url string) *routerGroup {
	g.Req = NewReq(Patch, joinPaths(g.basePath, url), g.out)
	return g
}

func (g *routerGroup) HEAD(url string) *routerGroup {
	g.Req = NewReq(Head, joinPaths(g.basePath, url), g.out)
	return g
}

func (g *routerGroup) OPTIONS(url string) *routerGroup {
	g.Req = NewReq(Options, joinPaths(g.basePath, url), g.out)
	return g
}

func (g *routerGroup) SetBody(obj interface{}) *routerGroup {
	g.Req.bodyEncoder = encode.NewBodyEncode(obj)
	return g
}

func (g *routerGroup) SetForm(obj interface{}) *routerGroup {
	g.Req.formEncode = obj
	return g
}

func (g *routerGroup) SetQuery(obj interface{}) *routerGroup {
	g.Req.queryEncode = obj
	return g
}

func (g *routerGroup) SetHeader(obj interface{}) *routerGroup {
	g.Req.headerEncode = obj
	return g
}

func (g *routerGroup) SetJSON(obj interface{}) *routerGroup {
	g.Req.bodyEncoder = encode.NewJSONEncode(obj)
	return g
}

func (g *routerGroup) SetXML(obj interface{}) *routerGroup {
	g.Req.bodyEncoder = encode.NewXMLEncode(obj)
	return g
}

func (g *routerGroup) SetYAML(obj interface{}) *routerGroup {
	g.Req.bodyEncoder = encode.NewYAMLEncode(obj)
	return g
}

// 此函数会修改Transport
func (g *routerGroup) UnixSocket(path string) *routerGroup {
	if g.out.Client.Transport == nil {
		g.out.Client.Transport = &http.Transport{}
	}

	transport, ok := g.out.Client.Transport.(*http.Transport)
	if !ok {
		g.Req.err = fmt.Errorf("UnixSocket:not found http.transport:%T", g.out.Client.Transport)
		return g
	}

	transport.Dial = func(proto, addr string) (conn net.Conn, err error) {
		return net.Dial("unix", path)
	}

	return g
}

// 此函数会修改Transport
func (g *routerGroup) SetProxy(proxyURL string) *routerGroup {
	proxy, err := url.Parse(modifyURL(proxyURL))
	if err != nil {
		g.Req.err = err
		return g
	}

	if g.out.Client.Transport == nil {
		g.out.Client.Transport = &http.Transport{}
	}

	transport, ok := g.out.Client.Transport.(*http.Transport)
	if !ok {
		g.Req.err = fmt.Errorf("SetProxy:not found http.transport:%T", g.out.Client.Transport)
		return g
	}

	transport.Proxy = http.ProxyURL(proxy)

	return g
}

func (g *routerGroup) SetCookies(c ...*http.Cookie) *routerGroup {
	g.Req.cookies = append(g.Req.cookies, c...)
	return g
}

func (g *routerGroup) BindBody(obj interface{}) *routerGroup {
	g.Req.bodyDecoder = decode.NewBodyDecode(obj)
	return g
}

func (g *routerGroup) BindHeader(obj interface{}) *routerGroup {
	g.Req.headerDecode = obj
	return g
}

func (g *routerGroup) BindJSON(obj interface{}) *routerGroup {
	g.Req.bodyDecoder = decode.NewJSONDecode(obj)
	return g
}

func (g *routerGroup) BindXML(obj interface{}) *routerGroup {
	g.Req.bodyDecoder = decode.NewXMLDecode(obj)
	return g
}

func (g *routerGroup) BindYAML(obj interface{}) *routerGroup {
	g.Req.bodyDecoder = decode.NewYAMLDecode(obj)
	return g
}

func (g *routerGroup) Code(httpCode *int) *routerGroup {
	g.Req.httpCode = httpCode
	return g
}

func (g *routerGroup) Callback(cb func(*Context) error) *routerGroup {
	g.Req.callback = cb
	return g
}

func (g *routerGroup) WithContext(c context.Context) *routerGroup {
	g.Req.c = c
	return g
}

func (g *routerGroup) Debug(d ...interface{}) *routerGroup {
	fmt.Printf("group option = %p\n", &g.Req.g.opt)
	for _, v := range d {
		switch opt := v.(type) {
		case bool:
			g.Req.g.opt.Debug = true
		case DebugOpt:
			opt.Apply(&g.Req.g.opt)
		}
	}

	return g
}

func (g *routerGroup) Do() (err error) {
	defer func() {
		g.Req = nil
		g.reqs = g.reqs[0:0]
	}()

	if g.Req != nil {
		if err = g.Req.Do(); err != nil {
			return err
		}
	}

	for _, r := range g.reqs {
		if err = r.Do(); err != nil {
			return err
		}
	}

	return nil
}

func (g *routerGroup) Next() *routerGroup {
	g.reqs = append(g.reqs, g.Req)
	g.Req = nil
	return g
}
