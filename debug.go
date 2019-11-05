package gout

import (
	"bytes"
	"fmt"
	"github.com/guonaihong/gout/color"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func ToBodyType(s string) color.BodyType {
	switch strings.ToLower(s) {
	case "json":
		return color.JsonType
	case "xml":
		return color.XmlType
	case "yaml":
		return color.YamlType
	}

	return color.TxtType
}

type DebugOption struct {
	Write       io.Writer
	Debug       bool
	Color       bool
	ReqBodyType string
	RspBodyType string
}

type DebugOpt interface {
	Apply(*DebugOption)
}

type DebugFunc func(*DebugOption)

func (f DebugFunc) Apply(o *DebugOption) {
	f(o)
}

func DebugColor() DebugOpt {
	return DebugFunc(func(o *DebugOption) {
		o.Color = true
		o.Debug = true
		o.Write = os.Stdout
	})
}

func (do *DebugOption) resetBodyAndPrint(req *http.Request, resp *http.Response) error {
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	resp.Body = ioutil.NopCloser(bytes.NewReader(all))
	err = do.debugPrint(req, resp)
	resp.Body = ioutil.NopCloser(bytes.NewReader(all))
	return err
}

func (do *DebugOption) debugPrint(req *http.Request, rsp *http.Response) error {
	if t := rsp.Header.Get("Content-Type"); len(t) != 0 &&
		strings.Index(t, "application/json") != -1 {
		do.RspBodyType = "json"
	}

	if do.Write == nil {
		do.Write = os.Stdout
	}

	var w io.Writer = do.Write

	cl := color.New(do.Color)
	path := "/"
	if len(req.URL.Path) > 0 {
		path = req.URL.RequestURI()
	}

	fmt.Fprintf(w, "> %s %s %s\r\n", req.Method, path, req.Proto)

	// write request header
	for k, v := range req.Header {
		fmt.Fprintf(w, "> %s: %s\r\n", cl.Sgrayf(k),
			cl.Sbluef(strings.Join(v, ",")))
	}

	fmt.Fprint(w, ">\r\n")
	fmt.Fprint(w, "\n")

	// write req body
	if req.GetBody != nil {
		b, err := req.GetBody()
		if err != nil {
			return err
		}

		var r = io.Reader(b)
		format := color.NewFormatEncoder(r, do.Color, ToBodyType(do.ReqBodyType))
		if format != nil {
			r = format
		}

		if _, err := io.Copy(w, r); err != nil {
			return err
		}
		fmt.Fprintf(w, "\r\n\r\n")
	}

	fmt.Fprintf(w, "< %s %s\r\n", rsp.Proto, rsp.Status)
	for k, v := range rsp.Header {
		fmt.Fprintf(w, "< %s: %s\r\n", cl.Sgrayf(k), cl.Sbluef(strings.Join(v, ",")))
	}

	fmt.Fprintf(w, "\r\n\r\n")
	// write rsp body
	var r = io.Reader(rsp.Body)
	format := color.NewFormatEncoder(r, do.Color, ToBodyType(do.RspBodyType))
	if format != nil {
		r = format
	}
	_, err := io.Copy(w, r)

	fmt.Fprintf(w, "\r\n\r\n")

	return err
}
