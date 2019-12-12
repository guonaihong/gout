package gout

import (
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
