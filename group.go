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
	Req
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
	g.Req = reqDef(Get, joinPaths(g.basePath, url), g.out)
	return g
}

func (g *routerGroup) POST(url string) *routerGroup {
	g.Req = reqDef(Post, joinPaths(g.basePath, url), g.out)
	return g
}

func (g *routerGroup) PUT(url string) *routerGroup {
	g.Req = reqDef(Put, joinPaths(g.basePath, url), g.out)
	return g
}

func (g *routerGroup) DELETE(url string) *routerGroup {
	g.Req = reqDef(Delete, joinPaths(g.basePath, url), g.out)
	return g
}

func (g *routerGroup) PATCH(url string) *routerGroup {
	g.Req = reqDef(Patch, joinPaths(g.basePath, url), g.out)
	return g
}

func (g *routerGroup) HEAD(url string) *routerGroup {
	g.Req = reqDef(Head, joinPaths(g.basePath, url), g.out)
	return g
}

func (g *routerGroup) OPTIONS(url string) *routerGroup {
	g.Req = reqDef(Options, joinPaths(g.basePath, url), g.out)
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

func (g *routerGroup) SetWWWForm(obj interface{}) *routerGroup {
	g.Req.bodyEncoder = encode.NewWWWFormEncode(obj)
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
	g.out.opt.ReqBodyType = "json"
	g.Req.bodyEncoder = encode.NewJSONEncode(obj)
	return g
}

func (g *routerGroup) SetXML(obj interface{}) *routerGroup {
	g.out.opt.ReqBodyType = "xml"
	g.Req.bodyEncoder = encode.NewXMLEncode(obj)
	return g
}

func (g *routerGroup) SetYAML(obj interface{}) *routerGroup {
	g.out.opt.ReqBodyType = "yaml"
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
	g.out.opt.RspBodyType = "json"
	g.Req.bodyDecoder = decode.NewJSONDecode(obj)
	return g
}

func (g *routerGroup) BindXML(obj interface{}) *routerGroup {
	g.out.opt.RspBodyType = "xml"
	g.Req.bodyDecoder = decode.NewXMLDecode(obj)
	return g
}

func (g *routerGroup) BindYAML(obj interface{}) *routerGroup {
	g.out.opt.RspBodyType = "yaml"
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
	for _, v := range d {
		switch opt := v.(type) {
		case bool:
			defaultDebug(&g.out.opt)
		case DebugOpt:
			opt.Apply(&g.out.opt)
		}
	}

	return g
}

func (g *routerGroup) Do() (err error) {
	defer func() {
		g.Req.Reset()
		g.reqs = g.reqs[0:0]
	}()

	for _, r := range g.reqs {
		if err = r.Do(); err != nil {
			return err
		}
	}

	if err = g.Req.Do(); err != nil {
		return err
	}

	return nil
}

func (g *routerGroup) Next() *routerGroup {
	r := g.Req.clone()
	g.Req.Reset()
	g.Req.url = ""
	g.reqs = append(g.reqs, &r)
	return g
}

func (g *routerGroup) FilterBench() *Bench {
	return &Bench{g: g}
}
