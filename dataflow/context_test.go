package dataflow

import (
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
	r := gin.Default()
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
	for _, p := range path {
		err := New().GET(ts.URL + "/" + p).Callback(func(c *Context) error {
			assert.NotEqual(t, c.Code, 404)

			switch c.Code {
			case 500:
				var s string
				err := c.BindBody(&s)
				assert.NoError(t, err)
				assert.Equal(t, "fail", s)

				count++
			case 200:
				type jsonResult struct {
					Errcode int    `json:"errcode"`
					Errmsg  string `json:"errmsg"`
				}

				var j jsonResult
				err := c.BindJSON(&j)
				assert.NoError(t, err)
				assert.Equal(t, j.Errmsg, "ok")
				assert.Equal(t, j.Errcode, 0)
				count++
			}

			var tHeader testHeader
			err := c.BindHeader(&tHeader)
			assert.NoError(t, err)

			assert.Equal(t, tHeader.HeaderKey, "headerVal")
			return nil
		}).Do()
		assert.NoError(t, err)

	}
	assert.Equal(t, count, 2)
}

func testServrBodyYAML(t *testing.T) *gin.Engine {
	r := gin.Default()
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
	for _, p := range path {
		err := New().GET(ts.URL + "/" + p).Callback(func(c *Context) error {
			assert.NotEqual(t, c.Code, 404)

			switch c.Code {
			case 500:
				var s string
				err := c.BindBody(&s)
				assert.NoError(t, err)
				assert.Equal(t, "fail", s)

				count++
			case 200:
				type yamlResult struct {
					Errcode int    `yaml:"errcode"`
					Errmsg  string `yaml:"errmsg"`
				}

				var y yamlResult
				err := c.BindYAML(&y)
				assert.NoError(t, err)
				assert.Equal(t, y.Errmsg, "ok")
				assert.Equal(t, y.Errcode, 0)
				count++
			}

			return nil
		}).Do()
		assert.NoError(t, err)

	}
	assert.Equal(t, count, 2)
}

func testServrBodyXML(t *testing.T) *gin.Engine {
	r := gin.Default()
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
	for _, p := range path {
		err := New().GET(ts.URL + "/" + p).Callback(func(c *Context) error {
			assert.NotEqual(t, c.Code, 404)

			switch c.Code {
			case 500:
				var s string
				err := c.BindBody(&s)
				assert.NoError(t, err)
				assert.Equal(t, "fail", s)

				count++
			case 200:
				type xmlResult struct {
					Errcode int    `xml:"errcode"`
					Errmsg  string `xml:"errmsg"`
				}

				var x xmlResult
				err := c.BindXML(&x)
				assert.NoError(t, err)
				assert.Equal(t, x.Errmsg, "ok")
				assert.Equal(t, x.Errcode, 0)
				count++
			}

			var tHeader testHeader
			err := c.BindHeader(&tHeader)
			assert.NoError(t, err)

			assert.Equal(t, tHeader.HeaderKey, "headerVal")
			return nil
		}).Do()
		assert.NoError(t, err)

	}
	assert.Equal(t, count, 2)
}

func TestContext_fail(t *testing.T) {
	s := testServrBodyJSON(t)
	ts := httptest.NewServer(http.HandlerFunc(s.ServeHTTP))

	var j, j2 testContextJSON

	// BindJSON和Callback只能选一个
	err := GET(ts.URL + "/200").BindJSON(&j).Callback(func(c *Context) error {
		return c.BindJSON(&j2)
	}).Do()

	assert.Error(t, err)
}
