package dataflow

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	api "github.com/guonaihong/gout/interface"
	"github.com/stretchr/testify/assert"
)

type demoRequestMiddler struct{}

func (d *demoRequestMiddler) ModifyRequest(req *http.Request) error {
	req.Body = ioutil.NopCloser(strings.NewReader("demo"))
	req.ContentLength = 4
	return nil
}

func demoRequest() api.RequestMiddler {
	return &demoRequestMiddler{}
}

func Test_RequestModify(t *testing.T) {
	ts := createGeneralEcho()
	s := ""
	err := New().POST(ts.URL).RequestUse(demoRequest()).SetBody("hello").BindBody(&s).Do()
	assert.NoError(t, err, fmt.Sprintf("test url:%s", ts.URL))
	assert.Equal(t, s, "demo")
}
