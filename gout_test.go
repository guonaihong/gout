package gout

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestNew(t *testing.T) {
	c := &http.Client{}
	tests := []*Gout{
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

func TestMethodFail(t *testing.T) {
	var total int32

	router := setupMethod(&total)

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	out := New(nil)
	err := out.GET(ts.URL + "/someGet").Next().
		POST(ts.URL + "/somePost").Next().
		PUT(ts.URL + "/somePut").Next().
		DELETE(ts.URL + "/someDelete").Next().
		PATCH(ts.URL + "/somePatch").Next().
		HEAD("127.0.0.1" /*不存在*/ + "/someHead").Next().
		OPTIONS(ts.URL + "/someOptions").Do()

	assert.Error(t, err, fmt.Sprintf("total = %d", total))

}

func TestMethod(t *testing.T) {

	var total int32

	router := setupMethod(&total)

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	out := New(nil)
	err := out.GET(ts.URL + "/someGet").Next().
		POST(ts.URL + "/somePost").Next().
		PUT(ts.URL + "/somePut").Next().
		DELETE(ts.URL + "/someDelete").Next().
		PATCH(ts.URL + "/somePatch").Next().
		HEAD(ts.URL + "/someHead").Next().
		OPTIONS(ts.URL + "/someOptions").Do()

	assert.NoError(t, err)

	assert.Equal(t, int(total), 7)

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
