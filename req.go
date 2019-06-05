package gout

import (
	"fmt"
	"net/http"
	"strings"
)

type Req struct {
	method string
	url    string

	g *Gout
}

func (r *Req) Do() (err error) {

	req, err := http.NewRequest(r.method, r.url, nil)
	if err != nil {
		return err
	}

	resp, err := r.g.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func modifyUrl(url string) string {
	if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://") {
		return url
	}

	if strings.HasPrefix(url, ":") {
		return fmt.Sprintf("http://127.0.0.1%s", url)
	}

	if strings.HasPrefix(url, "/") {
		return fmt.Sprintf("http://127.0.0.1%s", url)
	}

	return fmt.Sprintf("http://%s", url)
}

func NewReq(method string, url string, g *Gout) *Req {
	return &Req{method: method, url: modifyUrl(url), g: g}
}
