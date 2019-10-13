package gout

import (
	"bytes"
	"errors"
	"fmt"
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
			req.GetBody = func() (io.ReadCloser, error) { return &ReadCloseFail{}, errors.New("fail") }
			rsp.Body = ioutil.NopCloser(bytes.NewReader(nil))
			return req, &rsp
		},
		func() (*http.Request, *http.Response) {
			// GetBody不为空的情况, 但是io.Copy时候返回错误
			req, _ := http.NewRequest("GET", "/", nil)
			rsp := http.Response{}
			req.GetBody = func() (io.ReadCloser, error) { return &ReadCloseFail{}, nil }
			rsp.Body = ioutil.NopCloser(bytes.NewReader(nil))
			return req, &rsp
		},
		func() (*http.Request, *http.Response) {
			// rsp.Body出错的情况
			req, _ := http.NewRequest("GET", "/", &bytes.Buffer{})
			req.GetBody = func() (io.ReadCloser, error) { return ioutil.NopCloser(bytes.NewReader(nil)), nil }
			rsp := http.Response{}
			rsp.Body = ioutil.NopCloser(&ReadCloseFail{})
			return req, &rsp
		},
	}

	for i, c := range test {
		req, rsp := c()
		err := resetBodyAndPrint(req, rsp)
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
	for i, c := range test {
		req, rsp := c()
		err := resetBodyAndPrint(req, rsp)
		assert.NoError(t, err, fmt.Sprintf("test index = %d", i))
	}
}
