package dataflow

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
		switch tc.Method {
		case "get":
			New().GET("{{.Host}}/{{.Method}}", tc).BindBody(&body).Do()
		case "post":
			New().POST("{{.Host}}/{{.Method}}", tc).BindBody(&body).Do()
		case "put":
			New().PUT("{{.Host}}/{{.Method}}", tc).BindBody(&body).Do()
		case "patch":
			New().PATCH("{{.Host}}/{{.Method}}", tc).BindBody(&body).Do()
		case "options":
			New().OPTIONS("{{.Host}}/{{.Method}}", tc).BindBody(&body).Do()
		case "head":
			code := 0
			New().HEAD("{{.Host}}/{{.Method}}", tc).Debug(true).BindBody(&body).Code(&code).Do()
			New().SetMethod(strings.ToUpper(tc.Method)).SetURL("{{.Host}}/{{.Method}}", tc).Debug(true).BindBody(&body2).Code(&code).Do()
			assert.Equal(t, code, 200)
			continue
		}

		New().SetMethod(strings.ToUpper(tc.Method)).SetURL("{{.Host}}/{{.Method}}", tc).Debug(true).BindBody(&body2).Do()
		assert.Equal(t, body, tc.Method)
		b := assert.Equal(t, body2, tc.Method)
		if !b {
			return
		}
	}

}
