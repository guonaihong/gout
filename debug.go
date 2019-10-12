package gout

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func resetBodyAndPrint(req *http.Request, resp *http.Response) error {
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	resp.Body = ioutil.NopCloser(bytes.NewReader(all))
	err = debugPrint(req, resp)
	resp.Body = ioutil.NopCloser(bytes.NewReader(all))
	return err
}

func debugPrint(req *http.Request, rsp *http.Response) error {
	var w io.Writer = os.Stdout

	path := "/"
	if len(req.URL.Path) > 0 {
		path = req.URL.RequestURI()
	}

	fmt.Fprintf(w, "> %s %s %s\r\n", req.Method, path, req.Proto)

	// write request header
	for k, v := range req.Header {
		fmt.Fprintf(w, "> %s: %s\r\n", k, strings.Join(v, ","))
		fmt.Fprintf(w, "> %s: %s\r\n", k,
			strings.Join(v, ","))
	}

	fmt.Fprint(w, ">\r\n")
	fmt.Fprint(w, "\n")

	// write body
	if b, err := req.GetBody(); err != nil {
		return err
	} else {
		io.Copy(os.Stdout, b)
		fmt.Fprintf(w, "\r\n\r\n")
	}

	fmt.Fprintf(w, "< %s %s\r\n", rsp.Proto, rsp.Status)
	for k, v := range rsp.Header {
		fmt.Fprintf(w, "< %s: %s\r\n", k, strings.Join(v, ","))
		fmt.Fprintf(w, "< %s: %s\r\n", k,
			strings.Join(v, ","))
	}

	_, err := io.Copy(w, rsp.Body)
	return err
}
