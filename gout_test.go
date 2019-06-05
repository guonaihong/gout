package gout

import (
	"github.com/gin-gonic/gin"
	"sync/atomic"
	"testing"
	"time"
)

func TestMethod(t *testing.T) {

	var total int32

	router := func() {
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

		// By default it serves on :8080 unless a
		// PORT environment variable was defined.
		router.Run()
		// router.Run(":3000") for a hard coded port
	}

	go router()

	time.Sleep(time.Millisecond * 250)

	out := New(nil)
	err := out.GET(":8080/someGet").Next().
		POST(":8080/somePost").Next().
		PUT(":8080/somePut").Next().
		DELETE(":8080/someDelete").Next().
		PATCH(":8080/somePatch").Next().
		HEAD(":8080/someHead").Next().
		OPTIONS(":8080/someOptions").Next().Do()

	if err != nil {
		t.Errorf("http client fail:%v\n", err)
	}

	if total != 7 {
		t.Errorf("got %d want 7\n", total)
	}
}
