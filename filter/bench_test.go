package gout

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/guonaihong/gout/core"
	"github.com/guonaihong/gout/dataflow"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

const (
	benchNumber = 300
	benchTime   = 500 * time.Millisecond
)

func setupBenchNumber(total *int32) *gin.Engine {
	router := gin.New()
	router.POST("/", func(c *gin.Context) {
		atomic.AddInt32(total, 1)
		c.String(200, "12345")
	})

	return router
}

// 测试压测次数
func Test_Bench_Number(t *testing.T) {
	total := int32(0)
	router := setupBenchNumber(&total)
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

	err := dataflow.POST(ts.URL).
		SetJSON(core.H{"hello": "world"}).
		Filter().
		Bench().
		Concurrent(20).
		Number(benchNumber).
		Do()

	assert.Equal(t, total, int32(benchNumber))
	assert.NoError(t, err)
}

// 测试压测时间
func Test_Bench_Durations(t *testing.T) {
	total := int32(0)
	router := setupBenchNumber(&total)
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

	s := time.Now()
	err := dataflow.POST(ts.URL).
		SetJSON(core.H{"hello": "world"}).
		Filter().
		Bench().
		Concurrent(20).
		Durations(benchTime).
		Do()

	take := time.Now().Sub(s)

	assert.NoError(t, err)
	assert.LessOrEqual(t, int64(benchNumber-100*time.Millisecond), int64(take))
}

// 测试压测频率
func Test_Bench_Rate(t *testing.T) {
	total := int32(0)
	router := setupBenchNumber(&total)
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

	number := 800
	rate := 400
	s := time.Now()
	err := dataflow.POST(ts.URL).
		SetJSON(core.H{"hello": "world"}).
		Filter().
		Bench().
		Rate(rate).
		Concurrent(20).
		Number(number).
		Do()

	take := time.Now().Sub(s)

	assert.NoError(t, err)

	assert.LessOrEqual(t, int64(take), int64(time.Duration(time.Duration(number/rate)*time.Second+100*time.Millisecond)))
	assert.GreaterOrEqual(t, int64(take), int64(time.Duration(number/rate)*time.Second-time.Second))
}

// 测试自定义函数
func Test_Bench_Loop(t *testing.T) {
	uid := uuid.New()
	i := int32(0)

	total := int32(0)
	router := setupBenchNumber(&total)
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

	err := NewBench().
		Concurrent(25).
		Number(2000).
		Loop(func(c *dataflow.Context) error {
			id := atomic.AddInt32(&i, 1)
			c.POST(ts.URL).Debug(true).SetJSON(core.H{"sid": uid.String(),
				"appkey": fmt.Sprintf("ak:%d", id),
				"text":   fmt.Sprintf("test text :%d", id)})
			return nil

		}).Do()

	assert.NoError(t, err)
}
