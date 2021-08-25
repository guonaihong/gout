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

type demoResponseMiddler struct{}

func (d *demoResponseMiddler) ModifyResponse(response *http.Response) error {
	b := response.StatusCode == 200
	fmt.Println("是否200 ==", b)

	// 修改responseBody。 因为返回值大概率会有 { code, data,msg} 等字段。 这里想验证code. 如果不对就返回error。 对的话将 data中的内容写入body,这样后面BindJson的时候就可以直接处理业务了

	response.Body = ioutil.NopCloser(response.Body)

	return nil

}

func Test_ResponseModify(t *testing.T) {
	ts := createGeneralEcho()
	s := ""
	New().POST(ts.URL).SetBody("hello").ResponseUse(&demoResponseMiddler{}).BindBody(&s).Do()

}
