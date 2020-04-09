package dataflow

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/guonaihong/gout/color"
	"github.com/guonaihong/gout/core"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

// 测试resetBodyAndPrint出错
func TestResetBodyAndPrintFail(t *testing.T) {

	test := []func() (*http.Request, *http.Response){
		func() (*http.Request, *http.Response) {
			// GetBody为空的情况
			req, _ := http.NewRequest("GET", "/", nil)
			rsp := http.Response{}
			req.GetBody = func() (io.ReadCloser, error) { return nil, errors.New("fail") }
			rsp.Body = ioutil.NopCloser(bytes.NewReader(nil))
			return req, &rsp
		},
		func() (*http.Request, *http.Response) {
			// GetBody不为空的情况, 但是GetBody第二参数返回错误
			req, _ := http.NewRequest("GET", "/", nil)
			rsp := http.Response{}
			req.GetBody = func() (io.ReadCloser, error) { return &core.ReadCloseFail{}, errors.New("fail") }
			rsp.Body = ioutil.NopCloser(bytes.NewReader(nil))
			return req, &rsp
		},
		func() (*http.Request, *http.Response) {
			// GetBody不为空的情况, 但是io.Copy时候返回错误
			req, _ := http.NewRequest("GET", "/", nil)
			rsp := http.Response{}
			req.GetBody = func() (io.ReadCloser, error) { return &core.ReadCloseFail{}, nil }
			rsp.Body = ioutil.NopCloser(bytes.NewReader(nil))
			return req, &rsp
		},
		func() (*http.Request, *http.Response) {
			// rsp.Body出错的情况
			req, _ := http.NewRequest("GET", "/", &bytes.Buffer{})
			req.GetBody = func() (io.ReadCloser, error) { return ioutil.NopCloser(bytes.NewReader(nil)), nil }
			rsp := http.Response{}
			rsp.Body = ioutil.NopCloser(&core.ReadCloseFail{})
			return req, &rsp
		},
	}

	do := DebugOption{}
	for i, c := range test {
		req, rsp := c()
		err := do.resetBodyAndPrint(req, rsp)
		assert.Error(t, err, fmt.Sprintf("test index = %d", i))
	}
}

func TestResetBodyAndPrint(t *testing.T) {

	test := []func() (*http.Request, *http.Response){
		func() (*http.Request, *http.Response) {
			req, _ := http.NewRequest("GET", "/", nil)
			req.Header.Add("reqHeader1", "reqHeaderValue1")
			req.Header.Add("reqHeader2", "reqHeaderValue2")
			req.GetBody = func() (io.ReadCloser, error) { return ioutil.NopCloser(bytes.NewReader(nil)), nil }

			rsp := http.Response{}
			rsp.Body = ioutil.NopCloser(bytes.NewReader(nil))
			rsp.Header = http.Header{}
			rsp.Header.Add("rspHeader1", "rspHeaderValue1")
			rsp.Header.Add("rspHeader2", "rspHeaderValue2")
			return req, &rsp
		},
	}

	do := DebugOption{}
	for i, c := range test {
		req, rsp := c()
		err := do.resetBodyAndPrint(req, rsp)
		assert.NoError(t, err, fmt.Sprintf("test index = %d", i))
	}
}

func TestDebug_ToBodyType(t *testing.T) {
	type bodyTest struct {
		in   string
		need color.BodyType
	}

	tests := []bodyTest{
		{"json", color.JSONType},
		{"", color.TxtType},
	}

	for _, test := range tests {
		got := ToBodyType(test.in)
		assert.Equal(t, test.need, got)
	}
}

// 测试Debug接口在各个数据类型下面是否可以打印request和response消息
func TestDebug_Debug(t *testing.T) {

	rspVal := "hello world"
	ts := createGeneral(rspVal)
	defer ts.Close()

	var buf bytes.Buffer

	dbug := func() DebugOpt {
		return DebugFunc(func(o *DebugOption) {
			o.Debug = true
			o.Color = true
			o.Write = &buf
		})
	}()

	for index, err := range []error{
		// formdata
		func() error {
			buf.Reset()
			key := "testFormdata"
			val := "testFormdataValue"
			err := New().POST(ts.URL).Debug(dbug).SetForm(core.H{key: val}).Do()
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(key)), -1, core.BytesToString(buf.Bytes()))
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(val)), -1)
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(rspVal)), -1)
			return err
		}(),
		// json
		func() error {
			buf.Reset()
			key := "testkeyjson"
			val := "testvaluejson"
			err := New().POST(ts.URL).SetJSON(core.H{key: val}).Debug(dbug).Do()

			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(key)), -1, core.BytesToString(buf.Bytes()))
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(val)), -1)
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(rspVal)), -1)
			return err
		}(),
		// body
		func() error {
			buf.Reset()
			key := "testFormdata"
			err := New().POST(ts.URL).Debug(dbug).SetBody(key).Do()
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(key)), -1, core.BytesToString(buf.Bytes()))
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(rspVal)), -1)
			return err
		}(),
		// yaml
		func() error {
			buf.Reset()
			key := "testkeyyaml"
			val := "testvalueyaml"
			err := New().POST(ts.URL).Debug(dbug).SetYAML(core.H{key: val}).Do()
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(key)), -1, core.BytesToString(buf.Bytes()))
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(val)), -1)
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(rspVal)), -1)
			return err
		}(),
		// xml
		func() error {
			val := "testXMLValue"

			var d data
			d.Data = val
			err := New().POST(ts.URL).Debug(dbug).SetXML(d).Do()
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(val)), -1, core.BytesToString(buf.Bytes()))
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(rspVal)), -1)
			buf.Reset()
			return err
		}(),
		// x-www-form-urlencoded
		func() error {
			buf.Reset()
			key := "testwwwform"
			val := "testwwwformvalue"
			err := New().POST(ts.URL).Debug(dbug).SetWWWForm(core.H{key: val}).Do()
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(key)), -1, core.BytesToString(buf.Bytes()))
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(val)), -1)
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(rspVal)), -1)
			return err
		}(),
	} {
		assert.NoError(t, err, fmt.Sprintf("test index :%d", index))
		if err != nil {
			break
		}

	}
}
