package dataflow

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout/core"
	"github.com/stretchr/testify/assert"
)

func createBodyNil(accept *bool) *httptest.Server {
	r := gin.New()
	r.POST("/", func(c *gin.Context) {
		*accept = true
	})

	ts := httptest.NewServer(http.HandlerFunc(r.ServeHTTP))
	return ts
}

// 测试SetHeader SetQuery 混合
// SetBody/SetJSON/SetYAML 参数传递空指针
func TestSetXXX_Nil(t *testing.T) {
	accept := false

	for index, err := range []error{

		func() error {
			accept = false
			ts := createBodyNil(&accept)
			defer ts.Close()
			return New().GET(ts.URL).SetHeader(nil).SetQuery(nil).SetBody(nil).Do()
		}(),
		func() error {
			accept = false
			ts := createBodyNil(&accept)
			defer ts.Close()
			return New().GET(ts.URL).SetHeader(nil).SetQuery(nil).SetJSON(nil).Do()
		}(),
		func() error {
			accept = false
			ts := createBodyNil(&accept)
			defer ts.Close()
			return New().GET(ts.URL).SetHeader(nil).SetQuery(nil).SetXML(nil).Do()
		}(),
		func() error {
			accept = false
			ts := createBodyNil(&accept)
			defer ts.Close()
			return New().GET(ts.URL).SetHeader(nil).SetQuery(nil).SetYAML(nil).Do()
		}(),
		func() error {
			accept = false
			ts := createBodyNil(&accept)
			defer ts.Close()
			return New().GET(ts.URL).SetHeader(nil).SetQuery(nil).SetProtoBuf(nil).Do()
		}(),
		func() error {
			accept = false
			ts := createBodyNil(&accept)
			defer ts.Close()
			return New().GET(ts.URL).SetHeader(nil).SetQuery(nil).SetWWWForm(nil).Do()
		}(),
		func() error {
			accept = false
			ts := createBodyNil(&accept)
			defer ts.Close()
			return New().GET(ts.URL).SetHeader(nil).SetQuery(nil).SetForm(nil).Do()
		}(),
	} {

		assert.NoError(t, err, fmt.Sprintf("test case index:%d", index))
		assert.False(t, accept)
	}

}

// 测试BindXXX参数传递空指针
func TestBindXXX_Nil(t *testing.T) {
	accept := false

	for index, err := range []error{

		func() error {
			accept = false
			ts := createBodyNil(&accept)
			defer ts.Close()
			return New().GET(ts.URL).SetHeader(nil).SetQuery(nil).BindBody(nil).Do()
		}(),
		func() error {
			accept = false
			ts := createBodyNil(&accept)
			defer ts.Close()
			return New().GET(ts.URL).SetHeader(nil).SetQuery(nil).BindJSON(nil).Do()
		}(),
		func() error {
			accept = false
			ts := createBodyNil(&accept)
			defer ts.Close()
			return New().GET(ts.URL).SetHeader(nil).SetQuery(nil).BindXML(nil).Do()
		}(),
		func() error {
			accept = false
			ts := createBodyNil(&accept)
			defer ts.Close()
			return New().GET(ts.URL).SetHeader(nil).SetQuery(nil).BindYAML(nil).Do()
		}(),
	} {

		assert.NoError(t, err, fmt.Sprintf("test case index:%d", index))
		assert.False(t, accept)
	}
}

func TestSetBody(t *testing.T) {

	router := func() *gin.Engine {
		router := gin.New()
		router.POST("/", func(c *gin.Context) {

			testBody := testBodyNeed{}

			assert.NoError(t, c.ShouldBindQuery(&testBody))

			var s string
			b := bytes.NewBuffer(nil)
			_, err := io.Copy(b, c.Request.Body)
			assert.NoError(t, err)
			defer c.Request.Body.Close()

			s = b.String()
			switch {
			case testBody.Int:
				assert.Equal(t, s, "1")
			case testBody.Int8:
				assert.Equal(t, s, "2")
			case testBody.Int16:
				assert.Equal(t, s, "3")
			case testBody.Int32:
				assert.Equal(t, s, "4")
			case testBody.Int64:
				assert.Equal(t, s, "5")
			case testBody.Uint:
				assert.Equal(t, s, "6")
			case testBody.Uint8:
				assert.Equal(t, s, "7")
			case testBody.Uint16:
				assert.Equal(t, s, "8")
			case testBody.Uint32:
				assert.Equal(t, s, "9")
			case testBody.Uint64:
				assert.Equal(t, s, "10")
			case testBody.String:
				assert.Equal(t, s, "test string")
			case testBody.Bytes:
				assert.Equal(t, s, "test bytes")
			case testBody.Float32:
				assert.Equal(t, s, "11")
			case testBody.Float64:
				assert.Equal(t, s, "12")
			case testBody.Reader:
				assert.Equal(t, s, "test io.Reader")
			default:
				c.JSON(500, "unknown type")
			}

		})

		return router
	}()

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

	code := 0
	err := New(nil).POST(ts.URL).SetQuery(core.H{"int": true}).SetBody(1).Code(&code).Do()
	assert.NoError(t, err)
	assert.Equal(t, code, 200)

	err = New(nil).POST(ts.URL).SetQuery(core.H{"int8": true}).SetBody(int8(2)).Code(&code).Do()
	assert.NoError(t, err)
	assert.Equal(t, code, 200)

	err = New(nil).POST(ts.URL).SetQuery(core.H{"int16": true}).SetBody(int16(3)).Code(&code).Do()
	assert.NoError(t, err)
	assert.Equal(t, code, 200)

	err = New(nil).POST(ts.URL).SetQuery(core.H{"int32": true}).SetBody(int32(4)).Code(&code).Do()
	assert.NoError(t, err)
	assert.Equal(t, code, 200)

	err = New(nil).POST(ts.URL).SetQuery(core.H{"int64": true}).SetBody(int64(5)).Code(&code).Do()
	assert.NoError(t, err)
	assert.Equal(t, code, 200)
	//=====================uint start

	err = New(nil).POST(ts.URL).SetQuery(core.H{"uint": true}).SetBody(6).Code(&code).Do()
	assert.NoError(t, err)
	assert.Equal(t, code, 200)

	err = New(nil).POST(ts.URL).SetQuery(core.H{"uint8": true}).SetBody(uint8(7)).Code(&code).Do()
	assert.NoError(t, err)
	assert.Equal(t, code, 200)

	err = New(nil).POST(ts.URL).SetQuery(core.H{"uint16": true}).SetBody(uint16(8)).Code(&code).Do()
	assert.NoError(t, err)
	assert.Equal(t, code, 200)

	err = New(nil).POST(ts.URL).SetQuery(core.H{"uint32": true}).SetBody(uint32(9)).Code(&code).Do()
	assert.NoError(t, err)
	assert.Equal(t, code, 200)

	err = New(nil).POST(ts.URL).SetQuery(core.H{"uint64": true}).SetBody(uint64(10)).Code(&code).Do()
	assert.NoError(t, err)
	assert.Equal(t, code, 200)
	//============================== float start

	err = New(nil).POST(ts.URL).SetQuery(core.H{"float32": true}).SetBody(float32(11)).Code(&code).Do()
	assert.NoError(t, err)
	assert.Equal(t, code, 200)

	err = New(nil).POST(ts.URL).SetQuery(core.H{"float64": true}).SetBody(float64(12)).Code(&code).Do()
	assert.NoError(t, err)
	assert.Equal(t, code, 200)

	err = New(nil).POST(ts.URL).SetQuery(core.H{"string": true}).SetBody("test string").Code(&code).Do()
	assert.NoError(t, err)
	assert.Equal(t, code, 200)

	// test bytes string
	err = New(nil).POST(ts.URL).SetQuery(core.H{"bytes": true}).SetBody([]byte("test bytes")).Code(&code).Do()
	assert.NoError(t, err)
	assert.Equal(t, code, 200)

	// test io.Reader
	err = New(nil).POST(ts.URL).SetQuery(core.H{"reader": true}).SetBody(bytes.NewBufferString("test io.Reader")).Code(&code).Do()
	assert.NoError(t, err)
	assert.Equal(t, code, 200)
}
func TestBindBody(t *testing.T) {
	router := func() *gin.Engine {
		router := gin.New()

		bodyBind := testBodyBind{}

		router.GET("/:type", func(c *gin.Context) {
			assert.NoError(t, c.ShouldBindUri(&bodyBind))

			switch bodyBind.Type {
			case "uint":
				c.String(200, "1")
			case "uint8":
				c.String(200, "2")
			case "uint16":
				c.String(200, "3")
			case "uint32":
				c.String(200, "4")
			case "uint64":
				c.String(200, "5")
			case "int":
				c.String(200, "6")
			case "int8":
				c.String(200, "7")
			case "int16":
				c.String(200, "8")
			case "int32":
				c.String(200, "9")
			case "int64":
				c.String(200, "10")
			case "float32":
				c.String(200, "11")
			case "float64":
				c.String(200, "12")
			case "string":
				c.String(200, "string")
			case "bytes":
				c.String(200, "bytes")
			case "io.writer":
				c.String(200, "io.writer")
			default:
				c.String(500, "unknown")
			}
		})

		return router
	}()

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

	tests := []testBodyReq{
		{url: "/uint", got: new(uint), need: core.NewPtrVal(uint(1))},
		{url: "/uint8", got: new(uint8), need: core.NewPtrVal(uint8(2))},
		{url: "/uint16", got: new(uint16), need: core.NewPtrVal(uint16(3))},
		{url: "/uint32", got: new(uint32), need: core.NewPtrVal(uint32(4))},
		{url: "/uint64", got: new(uint64), need: core.NewPtrVal(uint64(5))},
		{url: "/int", got: new(int), need: core.NewPtrVal(int(6))},
		{url: "/int8", got: new(int8), need: core.NewPtrVal(int8(7))},
		{url: "/int16", got: new(int16), need: core.NewPtrVal(int16(8))},
		{url: "/int32", got: new(int32), need: core.NewPtrVal(int32(9))},
		{url: "/int64", got: new(int64), need: core.NewPtrVal(int64(10))},
		{url: "/float32", got: new(float32), need: core.NewPtrVal(float32(11))},
		{url: "/float64", got: new(float64), need: core.NewPtrVal(float64(12))},
		{url: "/string", got: new(string), need: core.NewPtrVal("string")},
		{url: "/bytes", got: new([]byte), need: core.NewPtrVal([]byte("bytes"))},
		{url: "/io.writer", got: bytes.NewBufferString(""), need: bytes.NewBufferString("io.writer")},
	}

	for _, v := range tests {

		code := 0
		err := New(nil).GET(ts.URL + v.url).BindBody(v.got).Code(&code).Do()
		assert.Equal(t, code, 200)
		assert.NoError(t, err)
		assert.Equal(t, v.got, v.need)
	}

}
