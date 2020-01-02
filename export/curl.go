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

func (c *curl) formData(req *http.Request) error {
	contentType := req.Header.Get("Content-Type")
	if strings.Index(contentType, "multipart/form-data") == -1 {
		return nil
	}

	c.Data = ""

	//TODO 解析formdata
	return nil
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

		header[index] = fmt.Sprintf(`%s:%s`, hKey, strings.Join(hVal, ","))
		header[index] = fmt.Sprintf("%q", header[index])
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

	c.URL = fmt.Sprintf(`%q`, req.URL.String())
	c.Method = req.Method
	c.Data = fmt.Sprintf(`%q`, core.BytesToString(all))

	c.formData(req)
	tp := newTemplate(long)
	return tp.Execute(w, c)
}
