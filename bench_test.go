package gout

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

const bench_number = 100

func setup_bench_number(total *int32) *gin.Engine {
	router := gin.New()
	router.POST("/", func(c *gin.Context) {
		atomic.AddInt32(total, 1)
	})

	return router
}

func Test_Bench_Number(t *testing.T) {
	total := int32(0)
	router := setup_bench_number(&total)
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

	err := GET(ts.URL).
		SetJSON(H{"key": "val"}).
		FilterBench().
		Concurrent(20).
		Number(bench_number).
		Do()

	//assert.Equal(t, total, bench_number)//TODO open
	assert.NoError(t, err)
}
