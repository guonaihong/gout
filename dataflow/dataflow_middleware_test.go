package dataflow

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/guonaihong/gout/json"

	core "github.com/guonaihong/gout/core"

	"github.com/guonaihong/gout/middler"
	"github.com/stretchr/testify/assert"
)

type demoRequestMiddler struct{}

func (d *demoRequestMiddler) ModifyRequest(req *http.Request) error {
	req.Body = ioutil.NopCloser(strings.NewReader("demo"))
	req.ContentLength = 4
	return nil
}

func demoRequest() middler.RequestMiddler {
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
	obj := make(map[string]interface{})

	err = json.Unmarshal(all, &obj)
	if err != nil {
		return err
	}
	code := obj["code"]
	msg := obj["msg"]
	data := obj["data"]

	// Go中json中的数字经过反序列化会成为float64类型
	if float64(200) != code {
		return fmt.Errorf("请求失败, code %d msg %s", code, msg)
	} else {
		byt, _ := json.Marshal(&data)
		response.Body = ioutil.NopCloser(bytes.NewReader(byt))
		return nil
	}
}
func demoResponse() middler.ResponseMiddler {
	return &demoResponseMiddler{}
}

// 请求示例
func Test_ResponseModify(t *testing.T) {
	ts := createGeneralEcho()
	arrs := core.A{
		core.H{
			"code": 200,
			"msg":  "请求成功了",
			"data": core.H{
				"id":   "1",
				"name": "张三",
			},
		},
		core.H{
			"code": 500,
			"msg":  "查询数据库出错了",
			"data": nil,
		},
	}

	for i, arr := range arrs {
		// 返回值
		res := new(map[string]interface{})
		marshal, _ := json.Marshal(arr)

		err := New().POST(ts.URL).SetJSON(marshal).ResponseUse(demoResponse()).BindJSON(&res).Do()

		m := arr.(core.H)
		code := m["code"].(int)
		if code == 200 {
			assert.NoError(t, err, fmt.Sprintf("test case index:%d", i))
		} else {
			assert.Error(t, err, fmt.Sprintf("test case index:%d", i))
		}
		//log.Printf("请求 %d -->  参数 %s \n 响应 %s  \n  err  %s \n ", i, marshal, res, err)

	}
}
