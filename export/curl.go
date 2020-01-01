package export

import (
	"fmt"
	"github.com/guonaihong/gout/core"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

type Form struct {
	FileName    string
	ContentType string
	Data        string
}

type curl struct {
	Header   []string
	Method   string
	Data     string
	URL      string
	FormData []Form
}

func GenCurl(req *http.Request, long bool, w io.Writer) error {
	c := curl{}

	header := make([]string, 0, len(req.Header))

	for k := range req.Header {
		header = append(header, k)
	}

	sort.Strings(header)

	for index, hKey := range header {
		hVal := req.Header[hKey]

		// TODO转义
		header[index] = fmt.Sprintf("'%s:%s'", hKey, strings.Join(hVal, ","))
	}

	c.Header = header
	body, err := req.GetBody()
	if err != nil {
		return err
	}

	all, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	c.URL = fmt.Sprintf("'%s'", req.URL.String())
	c.Method = req.Method
	c.Data = fmt.Sprintf("'%s'", core.BytesToString(all))

	tp := newTemplate(long)
	return tp.Execute(w, c)
}
