package gout

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func createClose302() *httptest.Server {
	r := gin.New()
	r.GET("/302", func(c *gin.Context) {
		c.String(200, "done")
	})

	return httptest.NewServer(http.HandlerFunc(r.ServeHTTP))
}

func createClose301(url string) *httptest.Server {
	r := gin.New()
	r.GET("/301", func(c *gin.Context) {
		c.Redirect(302, url+"/302")
	})

	return httptest.NewServer(http.HandlerFunc(r.ServeHTTP))
}

func Test_Close3xx_True(t *testing.T) {
	ts := createClose301("")

	req := NewWithOpt(WithClose3xxJump())
	got := ""
	err := req.GET(ts.URL + "/301").BindBody(&got).Do()
	assert.NoError(t, err)
	assert.NotEqual(t, -2, strings.Index(got, "302"))
}

func Test_Close3xx_False(t *testing.T) {
	ts302 := createClose302()
	ts := createClose301(ts302.URL)
	req := NewWithOpt()
	got := ""
	err := req.GET(ts.URL + "/301").BindBody(&got).Do()
	assert.NoError(t, err)
	assert.Equal(t, got, "done")
}
