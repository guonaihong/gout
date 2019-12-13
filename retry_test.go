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

func setup_retry() *gin.Engine {
	router := gin.New()

	router.GET("/", func(c *gin.Context) {
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
	router := setup_retry()
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

	err := GET(ts.URL).
		Filter().
		Retry().
		Attempt(3).
		WaitTime(time.Millisecond * 10).
		MaxWaitTime(time.Millisecond * 800).
		Do()

	assert.NoError(t, err)
}
