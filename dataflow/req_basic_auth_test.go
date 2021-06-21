package dataflow

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func testServerBasicAuth(t *testing.T, userName, password string) *gin.Engine {
	r := gin.New()
	r.GET("/basicauth", func(c *gin.Context) {
		u, p, ok := c.Request.BasicAuth()
		if !ok {
			assert.Error(t, errors.New("basicauth get fail"))
			c.String(500, "fail")
			return
		}

		if userName != u {
			assert.Error(t, errors.New("user name fail"))
			c.String(500, "user name fail")
			return

		}

		if password != p {
			assert.Error(t, errors.New("password fail"))
			c.String(500, "password fail")
			return

		}

		c.String(200, "ok")
	})
	return r
}

func Test_BasicAuth(t *testing.T) {
	const (
		userName = "test-name"
		password = "test-passowrd"
	)

	s := testServerBasicAuth(t, userName, password)
	ts := httptest.NewServer(http.HandlerFunc(s.ServeHTTP))

	code := 0
	err := New().GET(ts.URL+"/basicauth").SetBasicAuth(userName, password).Code(&code).Do()
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
}
