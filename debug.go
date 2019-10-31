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

type DebugOption struct {
	Write io.Writer
	Debug bool
	Color bool
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

	// write body
	if req.GetBody != nil {
		b, err := req.GetBody()
		if err != nil {
			return err
		}

		if _, err := io.Copy(w, b); err != nil {
			return err
		}
		fmt.Fprintf(w, "\r\n\r\n")
	}

	fmt.Fprintf(w, "< %s %s\r\n", rsp.Proto, rsp.Status)
	for k, v := range rsp.Header {
		fmt.Fprintf(w, "< %s: %s\r\n", cl.Sgrayf(k), cl.Sbluef(strings.Join(v, ",")))
	}

	fmt.Fprintf(w, "\r\n\r\n")
	_, err := io.Copy(w, rsp.Body)

	fmt.Fprintf(w, "\r\n\r\n")

	return err
}
