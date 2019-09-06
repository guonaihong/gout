package gout

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

func testServr(t *testing.T) *gin.Engine {
	r := gin.Default()
	r.GET("/:code", func(c *gin.Context) {

		code := testContextStruct{}

		err := c.ShouldBindUri(&code)
		assert.NoError(t, err)

		switch code.Code {
		case 200:
			c.JSON(200, gin.H{"errmsg": "ok", "errcode": 0})
		case 500:
			c.String(500, "fail")
		}
	})

	return r
}

func TestContextBind(t *testing.T) {
	s := testServr(t)

	ts := httptest.NewServer(http.HandlerFunc(s.ServeHTTP))

	path := []string{"200", "500"}
	count := 0
	for _, p := range path {
		err := Def().GET(ts.URL + "/" + p).Callback(func(c *Context) error {
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

			return nil
		}).Do()
		assert.NoError(t, err)

	}
	assert.Equal(t, count, 2)
}
