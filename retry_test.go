package gout

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	retry_Count        = 3
	retry_doesNotExist = ":6364"
)

func setup_retry_fail() *gin.Engine {
	router := gin.New()

	var done chan struct{}
	router.GET("/", func(c *gin.Context) {
		<-done
	})

	return router
}

func setup_retry_ok() *gin.Engine {
	router := gin.New()

	router.GET("/", func(c *gin.Context) {
		c.String(200, "ok")
	})

	return router
}

func Test_Retry_min(t *testing.T) {
	type minData struct {
		a, b int
		need int
	}
	test := []minData{
		{a: 3, b: 4, need: 3},
		{a: 4, b: 3, need: 3},
		{a: 4, b: 4, need: 4},
	}

	r := Retry{}
	for _, v := range test {
		assert.Equal(t, r.min(time.Duration(v.a), time.Duration(v.b)), time.Duration(v.need))
	}
}

func Test_Retry_sleep(t *testing.T) {
	r := Retry{attempt: 100}
	r.init()

	// 方便画出曲线图
	for i := 0; i < r.attempt; i++ {
		sleep := r.getSleep()
		fmt.Printf("%d\n", sleep)
		//fmt.Printf("%d,%v\n", sleep, sleep)
		r.currAttempt++
	}
}

func Test_Retry_init(t *testing.T) {
	r := Retry{attempt: RetryAttempt, maxWaitTime: RetryMaxWaitTime, waitTime: RetryWaitTime}
	r1 := Retry{}
	r1.init()
	assert.Equal(t, r1, r)
}
func Test_Retry_Do(t *testing.T) {
	// 6364是随便写的一个端口，如果CI/CD这台机器上有这个端口，就要换个不存在的
	router := setup_retry_fail()
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	urls := []string{ts.URL, retry_doesNotExist}
	// 测试全部超时的情况
	for _, u := range urls {
		err := GET(u).
			SetTimeout(10 * time.Millisecond).
			Debug(true).
			Filter().
			Retry().
			Attempt(retry_Count).
			WaitTime(time.Millisecond * 10).
			MaxWaitTime(time.Millisecond * 50).
			Do()
		assert.Error(t, err)
	}

	// 测试正确的情况
	router = setup_retry_ok()
	ts = httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	urls = []string{ts.URL}
	for _, u := range urls {
		err := GET(u).
			SetTimeout(10 * time.Millisecond).
			Debug(true).
			Filter().
			Retry().
			Attempt(retry_Count).
			WaitTime(time.Millisecond * 10).
			MaxWaitTime(time.Millisecond * 50).
			Do()
		assert.NoError(t, err)
	}
}
