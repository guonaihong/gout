package gout

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type testURLTemplateCase struct {
	Host   string
	Method string
}

func createMethodEcho() *httptest.Server {
	router := func() *gin.Engine {
		router := gin.New()

		router.GET("/get", func(c *gin.Context) {
			c.String(200, "get")
		})

		router.POST("/post", func(c *gin.Context) {
			c.String(200, "post")
		})

		router.PUT("/put", func(c *gin.Context) {
			c.String(200, "put")
		})

		router.PATCH("/patch", func(c *gin.Context) {
			c.String(200, "patch")
		})

		router.OPTIONS("/options", func(c *gin.Context) {
			c.String(200, "options")
		})

		router.HEAD("/head", func(c *gin.Context) {
			c.String(200, "head")
		})

		return router
	}()

	return httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
}

func Test_URL_Template(t *testing.T) {
	ts := createMethodEcho()
	for _, tc := range []testURLTemplateCase{
		{Host: ts.URL, Method: "get"},
		{Host: ts.URL, Method: "post"},
		{Host: ts.URL, Method: "put"},
		{Host: ts.URL, Method: "patch"},
		{Host: ts.URL, Method: "options"},
		{Host: ts.URL, Method: "head"},
	} {
		body := ""
		body2 := ""
		code := 0
		switch tc.Method {
		case "get":
			err := GET("{{.Host}}/{{.Method}}", tc).Debug(true).BindBody(&body).Code(&code).Do()
			assert.NoError(t, err)
		case "post":
			err := POST("{{.Host}}/{{.Method}}", tc).BindBody(&body).Code(&code).Do()
			assert.NoError(t, err)
		case "put":
			err := PUT("{{.Host}}/{{.Method}}", tc).BindBody(&body).Code(&code).Do()
			assert.NoError(t, err)
		case "patch":
			err := PATCH("{{.Host}}/{{.Method}}", tc).BindBody(&body).Code(&code).Do()
			assert.NoError(t, err)
		case "options":
			err := OPTIONS("{{.Host}}/{{.Method}}", tc).BindBody(&body).Code(&code).Do()
			assert.NoError(t, err)
		case "head":
			code := 0
			err := HEAD("{{.Host}}/{{.Method}}", tc).Debug(true).BindBody(&body).Code(&code).Do()
			assert.NoError(t, err)
			err = New().SetMethod(strings.ToUpper(tc.Method)).SetURL("{{.Host}}/{{.Method}}", tc).Debug(true).BindBody(&body2).Code(&code).Do()
			assert.NoError(t, err)
			assert.Equal(t, code, 200)
			continue
		}
		assert.Equal(t, code, 200)

		err := New().SetMethod(strings.ToUpper(tc.Method)).SetURL("{{.Host}}/{{.Method}}", tc).Debug(true).BindBody(&body2).Do()
		assert.NoError(t, err)
		assert.Equal(t, body, tc.Method)
		b := assert.Equal(t, body2, tc.Method)
		if !b {
			return
		}
	}

}
