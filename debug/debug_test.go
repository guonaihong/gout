package debug

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout/color"
	"github.com/guonaihong/gout/core"
	"github.com/stretchr/testify/assert"
)

func createGeneralEcho() *httptest.Server {
	router := func() *gin.Engine {
		router := gin.New()

		router.POST("/", func(c *gin.Context) {
			_, err := io.Copy(c.Writer, c.Request.Body)
			if err != nil {
				fmt.Printf("createGeneralEcho fail:%v\n", err)
			}
		})

		return router
	}()

	return httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
}

func createGeneral(data string) *httptest.Server {
	router := func() *gin.Engine {
		router := gin.New()

		router.POST("/", func(c *gin.Context) {
			if len(data) > 0 {
				c.String(200, data)
			}
		})

		return router
	}()

	return httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
}

type data struct {
	ID   int    `json:"id" xml:"id"`
	Data string `json:"data" xml:"data"`
}

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
		err := do.ResetBodyAndPrint(req, rsp)
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
		err := do.ResetBodyAndPrint(req, rsp)
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
