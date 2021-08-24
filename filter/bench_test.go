package filter

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/guonaihong/gout/bench"
	"github.com/guonaihong/gout/core"
	"github.com/guonaihong/gout/dataflow"
	"github.com/stretchr/testify/assert"
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
	defer ts.Close()

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
	defer ts.Close()

	s := time.Now()
	err := dataflow.POST(ts.URL).
		SetJSON(core.H{"hello": "world"}).
		Filter().
		Bench().
		Concurrent(20).
		Durations(benchTime).
		Do()

	take := time.Since(s)

	assert.NoError(t, err)
	assert.LessOrEqual(t, int64(benchNumber-100*time.Millisecond), int64(take))
}

// 测试压测频率
func Test_Bench_Rate(t *testing.T) {
	total := int32(0)
	router := setupBenchNumber(&total)
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

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

	take := time.Since(s)

	assert.NoError(t, err)

	assert.Equal(t, int32(number), total)
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
	defer ts.Close()

	err := NewBench().
		Concurrent(25).
		Number(benchNumber).
		Loop(func(c *dataflow.Context) error {
			id := atomic.AddInt32(&i, 1)
			c.POST(ts.URL).Debug(true).SetJSON(core.H{"sid": uid.String(),
				"appkey": fmt.Sprintf("ak:%d", id),
				"text":   fmt.Sprintf("test text :%d", id)})
			return nil

		}).Do()

	assert.NoError(t, err)
	assert.Equal(t, total, int32(benchNumber))
}

func Test_Bench_fail(t *testing.T) {

	total := int32(0)
	router := setupBenchNumber(&total)
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	tests := []dataflow.Bencher{
		dataflow.POST(ts.URL).SetBody(time.Time{}).Filter().Bench().Concurrent(25).Number(benchNumber),
	}

	for _, v := range tests {
		err := v.Do()
		assert.Error(t, err)
	}

	testErr := errors.New("Test_Bench_fail")
	var r bench.Report
	err := NewBench().
		Concurrent(25).
		Number(benchNumber).
		GetReport(&r).
		Loop(func(c *dataflow.Context) error {
			return testErr
		}).Do()

	assert.NoError(t, err)
	v, ok := r.ErrMsg[testErr.Error()]
	assert.True(t, ok)
	assert.Equal(t, v, benchNumber)
}
