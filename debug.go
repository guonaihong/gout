package gout

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func debugPrint(req *http.Request, rsp *http.Response) {
	var w io.Writer = os.Stdout

	path := "/"
	if len(req.URL.Path) > 0 {
		path = req.URL.RequestURI()
	}

	fmt.Fprintf(w, "> %s %s %s\r\n", req.Method, path, req.Proto)
	for k, v := range req.Header {
		fmt.Fprintf(w, "> %s: %s\r\n", k, strings.Join(v, ","))
		fmt.Fprintf(w, "> %s: %s\r\n", k,
			strings.Join(v, ","))
	}

	fmt.Fprint(w, ">\r\n")
	fmt.Fprint(w, "\n")
	fmt.Fprintf(w, "< %s %s\r\n", rsp.Proto, rsp.Status)
	for k, v := range rsp.Header {
		fmt.Fprintf(w, "< %s: %s\r\n", k, strings.Join(v, ","))
		fmt.Fprintf(w, "< %s: %s\r\n", k,
			strings.Join(v, ","))
	}

	io.Copy(w, rsp.Body)
}
