package dataflow

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout/core"
	"github.com/stretchr/testify/assert"
)

//本文件主要测试Response接口
// 1.正确的情况
// 2.错误的情况

type testCase struct {
	json    bool
	reqData string
}

// response 接口测试mock server
func createResponseMock(t *testing.T) *httptest.Server {
	r := gin.New()
	r.POST("/", func(c *gin.Context) {
		_, err := io.Copy(c.Writer, c.Request.Body)
		assert.NoError(t, err)
	})

	ts := httptest.NewServer(http.HandlerFunc(r.ServeHTTP))
	return ts
}

// 测试正确的情况
func TestResponse_Ok(t *testing.T) {
	ts := createResponseMock(t)
	for _, tc := range []testCase{
		{json: true, reqData: `{"a":"b"}`},
	} {
		g := New().POST(ts.URL)
		if tc.json {
			rsp, err := g.SetJSON(tc.reqData).Response()
			assert.NoError(t, err)
			if err != nil {
				return
			}

			var out strings.Builder
			_, err = io.Copy(&out, rsp.Body)
			assert.NoError(t, err)
			if err != nil {
				return
			}

			assert.NoError(t, rsp.Body.Close())
			assert.Equal(t, tc.reqData, out.String())
		}
	}
}

// 测试错误的情况
func TestResponse_Fail(t *testing.T) {
	port := core.GetNoPortExists()

	rsp, err := New().POST(":" + port).Response()
	assert.Nil(t, rsp)
	assert.Error(t, err)
}
