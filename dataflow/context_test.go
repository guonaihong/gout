package dataflow

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testContextStruct struct {
	Code int `uri:"code"`
}

type testHeader struct {
	HeaderKey string `header:"headerKey"`
}

type testContextJSON struct {
	ErrMsg  string `json:"errmsg"`
	ErrCode int    `json:"errcode"`
}

func testServrBodyJSON(t *testing.T) *gin.Engine {
	r := gin.New()
	r.GET("/:code", func(c *gin.Context) {

		code := testContextStruct{}

		err := c.ShouldBindUri(&code)
		assert.NoError(t, err)

		c.Header("headerKey", "headerVal")
		switch code.Code {
		case 200:
			c.JSON(200, gin.H{"errmsg": "ok", "errcode": 0})
		case 500:
			c.String(500, "fail")
		}

	})

	return r
}

func TestContext_BindBodyJSON(t *testing.T) {
	s := testServrBodyJSON(t)

	ts := httptest.NewServer(http.HandlerFunc(s.ServeHTTP))

	path := []string{"200", "500"}
	count := 0
	type jsonResult struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}

	var j jsonResult
	var tHeader testHeader
	var str string
	for _, p := range path {
		code := 0
		err := New().GET(ts.URL + "/" + p).Debug(true).Callback(func(c *Context) error {
			assert.NotEqual(t, c.Code, 404)

			switch c.Code {
			case 500:
				c.BindBody(&str)
				count++
			case 200:
				c.BindJSON(&j)
				count++
			}

			code = c.Code

			c.BindHeader(&tHeader)

			return nil
		}).Do()

		assert.NoError(t, err)
		if code == 500 {
			assert.Equal(t, str, "fail")
		} else if code == 200 {
			assert.Equal(t, j.Errmsg, "ok")
			assert.Equal(t, j.Errcode, 0)
		}
		assert.Equal(t, tHeader.HeaderKey, "headerVal")

	}
	assert.Equal(t, count, 2)
}

func testServrBodyYAML(t *testing.T) *gin.Engine {
	r := gin.New()
	r.GET("/:code", func(c *gin.Context) {

		code := testContextStruct{}

		err := c.ShouldBindUri(&code)
		assert.NoError(t, err)

		c.Header("headerKey", "headerVal")
		switch code.Code {
		case 200:
			c.YAML(200, gin.H{"errmsg": "ok", "errcode": 0})
		case 500:
			c.String(500, "fail")
		}
	})

	return r
}

func TestContext_BindBodyYAML(t *testing.T) {
	s := testServrBodyYAML(t)

	ts := httptest.NewServer(http.HandlerFunc(s.ServeHTTP))

	path := []string{"200", "500"}
	count := 0
	type yamlResult struct {
		Errcode int    `yaml:"errcode"`
		Errmsg  string `yaml:"errmsg"`
	}

	var y yamlResult
	var str string
	for _, p := range path {
		code := 0
		err := New().GET(ts.URL + "/" + p).Callback(func(c *Context) error {
			assert.NotEqual(t, c.Code, 404)

			switch c.Code {
			case 500:
				c.BindBody(&str)
				count++
				code = c.Code
			case 200:
				count++
				c.BindYAML(&y)
				code = c.Code
			}

			return nil
		}).Do()
		assert.NoError(t, err)
		if code == 500 {
			assert.Equal(t, "fail", str)
		} else if code == 200 {
			assert.Equal(t, y.Errmsg, "ok")
			assert.Equal(t, y.Errcode, 0)
		}

	}
	assert.Equal(t, count, 2)
}

func testServrBodyXML(t *testing.T) *gin.Engine {
	r := gin.New()
	r.GET("/:code", func(c *gin.Context) {

		code := testContextStruct{}

		err := c.ShouldBindUri(&code)
		assert.NoError(t, err)

		c.Header("headerKey", "headerVal")
		switch code.Code {
		case 200:
			c.XML(200, gin.H{"errmsg": "ok", "errcode": 0})
		case 500:
			c.String(500, "fail")
		}
	})

	return r
}

func TestContext_BindBodyXML(t *testing.T) {
	s := testServrBodyXML(t)

	ts := httptest.NewServer(http.HandlerFunc(s.ServeHTTP))

	path := []string{"200", "500"}
	count := 0
	var str string
	var tHeader testHeader
	type xmlResult struct {
		Errcode int    `xml:"errcode"`
		Errmsg  string `xml:"errmsg"`
	}
	var x xmlResult

	for _, p := range path {
		code := 0
		err := New().GET(ts.URL + "/" + p).Debug(true).Callback(func(c *Context) error {
			assert.NotEqual(t, c.Code, 404)

			switch c.Code {
			case 500:
				c.BindBody(&str)
				count++
			case 200:

				c.BindXML(&x)
				count++
			}

			code = c.Code
			c.BindHeader(&tHeader)

			return nil
		}).Do()

		assert.NoError(t, err)
		if code == 500 {
			assert.Equal(t, "fail", str)
		} else if code == 200 {
			assert.Equal(t, x.Errmsg, "ok")
			assert.Equal(t, x.Errcode, 0)
		}
		assert.Equal(t, tHeader.HeaderKey, "headerVal")
	}
	assert.Equal(t, count, 2)
}

func TestContext_fail(t *testing.T) {
	s := testServrBodyJSON(t)
	ts := httptest.NewServer(http.HandlerFunc(s.ServeHTTP))

	//var j, j2 testContextJSON
	var j testContextJSON

	errs := []error{

		/*
			 //BindJSON和Callback只能选一个
					GET(ts.URL + "/200").BindJSON(&j).Callback(func(c *Context) error {
						c.BindJSON(&j2)
						return nil
					}).Do(),
		*/

		GET(ts.URL + "/200").BindJSON(&j).Callback(func(c *Context) error {
			return errors.New("fail")
		}).Do(),
	}

	for id, e := range errs {
		assert.Error(t, e, fmt.Sprintf("test id:%d\n", id))
	}
}
