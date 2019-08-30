package gout

import (
	"github.com/guonaihong/gout/decode"
	"github.com/guonaihong/gout/encode"
)

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
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
	g.Req = NewReq(GET, joinPaths(g.basePath, url), g.out)
	return g
}

func (g *routerGroup) POST(url string) *routerGroup {
	g.Req = NewReq(POST, joinPaths(g.basePath, url), g.out)
	return g
}

func (g *routerGroup) PUT(url string) *routerGroup {
	g.Req = NewReq(PUT, joinPaths(g.basePath, url), g.out)
	return g
}

func (g *routerGroup) DELETE(url string) *routerGroup {
	g.Req = NewReq(DELETE, joinPaths(g.basePath, url), g.out)
	return g
}

func (g *routerGroup) PATCH(url string) *routerGroup {
	g.Req = NewReq(PATCH, joinPaths(g.basePath, url), g.out)
	return g
}

func (g *routerGroup) HEAD(url string) *routerGroup {
	g.Req = NewReq(HEAD, joinPaths(g.basePath, url), g.out)
	return g
}

func (g *routerGroup) OPTIONS(url string) *routerGroup {
	g.Req = NewReq(OPTIONS, joinPaths(g.basePath, url), g.out)
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
	g.Req.bodyEncoder = encode.NewJsonEncode(obj)
	return g
}

func (g *routerGroup) SetXML(obj interface{}) *routerGroup {
	g.Req.bodyEncoder = encode.NewXmlEncode(obj)
	return g
}

func (g *routerGroup) SetYAML(obj interface{}) *routerGroup {
	g.Req.bodyEncoder = encode.NewYamlEncode(obj)
	return g
}

func (g *routerGroup) BindHeader(obj interface{}) *routerGroup {
	g.Req.headerDecode = obj
	return g
}

func (g *routerGroup) BindJSON(obj interface{}) *routerGroup {
	g.Req.bodyDecoder = decode.NewJsonDecode(obj)
	return g
}

func (g *routerGroup) BindXML(obj interface{}) *routerGroup {
	g.Req.bodyDecoder = decode.NewXmlDecode(obj)
	return g
}

func (g *routerGroup) BindYAML(obj interface{}) *routerGroup {
	g.Req.bodyDecoder = decode.NewYamlDecode(obj)
	return g
}

func (g *routerGroup) Code(httpCode *int) *routerGroup {
	g.Req.httpCode = httpCode
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
