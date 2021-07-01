package dataflow

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout/core"
	"github.com/stretchr/testify/assert"
)

type testWWWForm struct {
	Int     int     `form:"int" www-form:"int,omitempty"`
	Float64 float64 `form:"float64" www-form:"float64,omitempty"`
	String  string  `form:"string" www-form:"string,omitempty"`
}

func setupWWWForm(t *testing.T, need testWWWForm) *gin.Engine {
	r := gin.New()

	r.POST("/", func(c *gin.Context) {
		wf := testWWWForm{}

		err := c.ShouldBind(&wf)

		assert.NoError(t, err)
		//err := c.ShouldBind(&wf)
		assert.Equal(t, need, wf)
	})

	return r
}

func TestWWWForm(t *testing.T) {
	need := testWWWForm{
		Int:     3,
		Float64: 3.14,
		String:  "test-www-Form",
	}

	router := setupWWWForm(t, need)
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

	for _, v := range [][]interface{}{
		// 测试使用一个结构体的情况
		[]interface{}{
			testWWWForm{
				Int:     3,
				Float64: 3.14,
				String:  "test-www-Form",
			},
		},
		// 测试两个结构体，同一个数据
		[]interface{}{
			testWWWForm{
				Int: 3,
			},
			testWWWForm{
				Float64: 3.14,
				String:  "test-www-Form",
			},
		},
		// 测试两种类型混用
		[]interface{}{
			core.H{"int": 3},
			testWWWForm{
				Float64: 3.14,
				String:  "test-www-Form",
			},
		},
	} {
		err := POST(ts.URL).Debug(true).SetWWWForm(v...).Do()
		assert.NoError(t, err)
	}
}
