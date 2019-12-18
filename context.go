package gout

import (
	"github.com/guonaihong/gout/decode"
	"net/http"
)

// Context struct
type Context struct {
	Code int //http code
	Resp *http.Response
}

// BindHeader parse http header to obj variable.
// obj must be a pointer variable
func (c *Context) BindHeader(obj interface{}) error {
	return decode.Header.Decode(c.Resp, obj)
}

// BindBody parse the variables in http body to obj.
// obj must be a pointer variable
func (c *Context) BindBody(obj interface{}) error {
	return decode.Body(c.Resp.Body, obj)
}

// BindJSON parse the json string in http body to obj.
// obj must be a pointer variable
func (c *Context) BindJSON(obj interface{}) error {
	return decode.JSON(c.Resp.Body, obj)
}

// BindYAML parse the yaml string in http body to obj.
// obj must be a pointer variable
func (c *Context) BindYAML(obj interface{}) error {
	return decode.YAML(c.Resp.Body, obj)
}

// BindXML parse the xml string in http body to obj.
// obj must be a pointer variable
func (c *Context) BindXML(obj interface{}) error {
	return decode.XML(c.Resp.Body, obj)
}
