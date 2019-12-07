package bench

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func setup_report_server(total *int32) *gin.Engine {
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		atomic.AddInt32(total, 1)
	})

	return router
}

func runReport(p *Report, number int) {
	p.Init()
	work := make(chan struct{})
	go func() {
		for i := 0; i < number; i++ {
			work <- struct{}{}
		}
		close(work)
	}()
	quit := make(chan struct{})
	go func() {
		p.Process(work)
		close(quit)
	}()

	<-quit
	p.Cancel()
	p.WaitAll()

}

func newRequest(url string) (*http.Request, error) {
	b := bytes.NewBufferString("hello")
	return http.NewRequest("GET", url, b)
}

// 测试正常情况, 次数
func Test_Bench_Report_number(t *testing.T) {
	const number = 1000

	total := int32(0)
	router := setup_report_server(&total)
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

	ctx := context.Background()

	req, err := newRequest(ts.URL)
	assert.NoError(t, err)

	p := NewReport(ctx, 1, number, time.Duration(0), req, http.DefaultClient)

	runReport(p, number)

	assert.Equal(t, int32(p.CompleteRequest), int32(number))
}

// 测试正常情况, 时间
func Test_Bench_Report_duration(t *testing.T) {
	const number = 1000

	total := int32(0)
	router := setup_report_server(&total)
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

	ctx := context.Background()

	req, err := newRequest(ts.URL)
	assert.NoError(t, err)

	p := NewReport(ctx, 1, 0, 300*time.Millisecond, req, http.DefaultClient)

	runReport(p, number)

	assert.Equal(t, int32(p.CompleteRequest), int32(number))
}

// 测试异常情况
func Test_Bench_Report_fail(t *testing.T) {
	const number = 1000

	ctx := context.Background()

	req, err := newRequest("fail")
	assert.NoError(t, err)

	p := NewReport(ctx, 1, number, time.Duration(0), req, http.DefaultClient)

	runReport(p, number)

	assert.Equal(t, int32(p.Failed), int32(number))
}
