package dataflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

// response拦截器修改示例
type demoResponseMiddler struct{}

func (d *demoResponseMiddler) ModifyResponse(response *http.Response) error {
	// 修改responseBody。 因为返回值大概率会有 { code, data,msg} 等字段,希望进行统一处理
	//这里想验证code. 如果不对就返回error。 对的话将 data中的内容写入body,这样后面BindJson的时候就可以直接处理业务了
	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	obj := map[string]string{}
	err = json.Unmarshal(all, &obj)
	if err != nil {
		return err
	}
	obj["ccc"] = "aaaa"
	byt, _ := json.Marshal(&obj)
	response.Body = ioutil.NopCloser(bytes.NewReader(byt))
	return nil

}
func demoResponse() api.ResponseMiddler {
	return &demoResponseMiddler{}
}

type Res struct {
	A   string `json:"a"`
	B   string `json:"b"`
	Ccc string `json:"ccc"`
}

func Test_ResponseModify(t *testing.T) {
	ts := createGeneralEcho()
	mp := map[string]string{
		"a": "aa",
		"b": "bb",
	}
	res := map[string]string{}
	//res :=""
	New().POST(ts.URL).SetBody("hello").SetJSON(mp).ResponseUse(demoResponse()).BindJSON(&res).Do()
	log.Printf("日志 -->  请求body %s , 响应body  %s", mp, res)
}
