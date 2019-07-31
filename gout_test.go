package gout

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestMethod(t *testing.T) {

	var total int32

	router := func() *gin.Engine {
		// Creates a gin router with default middleware:
		// logger and recovery (crash-free) middleware
		router := gin.Default()

		cb := func(c *gin.Context) {
			atomic.AddInt32(&total, 1)
		}

		router.GET("/someGet", cb)
		router.POST("/somePost", cb)
		router.PUT("/somePut", cb)
		router.DELETE("/someDelete", cb)
		router.PATCH("/somePatch", cb)
		router.HEAD("/someHead", cb)
		router.OPTIONS("/someOptions", cb)

		return router
	}()

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	out := New(nil)
	err := out.GET(ts.URL + "/someGet").Next().
		POST(ts.URL + "/somePost").Next().
		PUT(ts.URL + "/somePut").Next().
		DELETE(ts.URL + "/someDelete").Next().
		PATCH(ts.URL + "/somePatch").Next().
		HEAD(ts.URL + "/someHead").Next().
		OPTIONS(ts.URL + "/someOptions").Next().Do()

	assert.NoError(t, err)

	assert.Equal(t, int(total), 7)

}
