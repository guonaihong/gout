package gout

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout/dataflow"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestDebug(t *testing.T) {
	server := func() *gin.Engine {
		r := gin.New()

		r.GET("/", func(c *gin.Context) {
			all, err := ioutil.ReadAll(c.Request.Body)

			assert.NoError(t, err)
			c.String(200, string(all))
		})

		return r
	}()

	ts := httptest.NewServer(http.HandlerFunc(server.ServeHTTP))
	test := []func() DebugOpt{
		// 没有颜色输出
		NoColor,
		Trace,
	}

	s := ""
	for k, v := range test {
		s = ""
		err := GET(ts.URL).
			Debug(v()).
			SetBody(fmt.Sprintf("%d test debug.", k)).
			BindBody(&s).
			Do()
		assert.NoError(t, err)

		assert.Equal(t, fmt.Sprintf("%d test debug.", k), s)
	}
}

func TestNew(t *testing.T) {
	c := &http.Client{}
	tests := []*dataflow.Gout{
		New(nil),
		New(),
		New(c),
	}

	for _, v := range tests {
		assert.NotNil(t, v)
	}
}

func setupMethod(total *int32) *gin.Engine {

	router := gin.Default()

	cb := func(c *gin.Context) {
		atomic.AddInt32(total, 1)
	}

	router.GET("/someGet", cb)
	router.POST("/somePost", cb)
	router.PUT("/somePut", cb)
	router.DELETE("/someDelete", cb)
	router.PATCH("/somePatch", cb)
	router.HEAD("/someHead", cb)
	router.OPTIONS("/someOptions", cb)

	return router
}

func TestTopMethod(t *testing.T) {
	var total int32

	router := setupMethod(&total)

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	err := GET(ts.URL + "/someGet").Do()
	assert.NoError(t, err)

	err = POST(ts.URL + "/somePost").Do()
	assert.NoError(t, err)

	err = PUT(ts.URL + "/somePut").Do()
	assert.NoError(t, err)

	err = DELETE(ts.URL + "/someDelete").Do()
	assert.NoError(t, err)

	err = PATCH(ts.URL + "/somePatch").Do()
	assert.NoError(t, err)

	err = HEAD(ts.URL + "/someHead").Do()
	assert.NoError(t, err)

	err = OPTIONS(ts.URL + "/someOptions").Do()
	assert.NoError(t, err)

	assert.Equal(t, int(total), 7)
}
