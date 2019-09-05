package gout

import (
	"github.com/guonaihong/gout/decode"
	"net/http"
)

type Context struct {
	Code int //http code
	Resp *http.Response
}

func (c *Context) BindBody(obj interface{}) error {
	return decode.DecodeBody(c.Resp.Body, obj)
}

func (c *Context) BindJSON(obj interface{}) error {
	return decode.DecodeJSON(c.Resp.Body, obj)
}

func (c *Context) BindYAML(obj interface{}) error {
	return decode.DecodeYAML(c.Resp.Body, obj)
}

func (c *Context) BindXML(obj interface{}) error {
	return decode.DecodeXML(c.Resp.Body, obj)
}

func (c *Context) BindHeader(obj interface{}) error {
	return decode.Header.Decode(c.Resp, obj)
}
