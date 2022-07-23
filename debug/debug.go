package debug

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/guonaihong/gout/color"
)

// ToBodyType Returns the http body type,
// which mainly affects color highlighting
func ToBodyType(s string) color.BodyType {
	switch strings.ToLower(s) {
	case "json":
		return color.JSONType
	case "xml":
		//return color.XmlType //TODO open
	case "yaml":
		//return color.YamlType //TODO open
	}

	return color.TxtType
}

// Option Debug mode core data structure
type Option struct {
	Write       io.Writer
	Debug       bool
	Color       bool
	Trace       bool
	ReqBodyType string
	RspBodyType string
	TraceInfo
}

// DebugOpt is an interface for operating Option
type DebugOpt interface {
	Apply(*Option)
}

// DebugFunc DebugOpt is a function that manipulates core data structures
type DebugFunc func(*Option)

// Apply is an interface for operating Option
func (f DebugFunc) Apply(o *Option) {
	f(o)
}

func DefaultDebug(o *Option) {
	o.Color = true
	o.Debug = true
	o.Write = os.Stdout
}

// NoColor Turn off color highlight debug mode
func NoColor() DebugOpt {
	return DebugFunc(func(o *Option) {
		o.Color = false
		o.Debug = true
		o.Trace = true
		o.Write = os.Stdout
	})
}

func (do *Option) ResetBodyAndPrint(req *http.Request, resp *http.Response) error {
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	resp.Body.Close()

	resp.Body = ioutil.NopCloser(bytes.NewReader(all))
	err = do.debugPrint(req, resp)
	resp.Body = ioutil.NopCloser(bytes.NewReader(all))
	return err
}

func (do *Option) debugPrint(req *http.Request, rsp *http.Response) error {
	if t := req.Header.Get("Content-Type"); len(t) > 0 && strings.Contains(t, "json") {
		do.ReqBodyType = "json"
	}

	if t := rsp.Header.Get("Content-Type"); len(t) != 0 &&
		strings.Contains(t, "application/json") {
		do.RspBodyType = "json"
	}

	if do.Write == nil {
		do.Write = os.Stdout
	}

	var w io.Writer = do.Write

	cl := color.New(do.Color)
	path := "/"
	if len(req.URL.RequestURI()) > 0 {
		path = req.URL.RequestURI()
	}

	fmt.Fprintf(w, "> %s %s %s\r\n", req.Method, path, req.Proto)

	// write request header
	for k, v := range req.Header {
		fmt.Fprintf(w, "> %s: %s\r\n", cl.Spurple(k),
			cl.Sblue(strings.Join(v, ",")))
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
		fmt.Fprintf(w, "< %s: %s\r\n", cl.Spurple(k), cl.Sblue(strings.Join(v, ",")))
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
