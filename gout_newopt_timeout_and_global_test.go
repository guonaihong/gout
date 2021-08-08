package gout

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupDataFlow(t *testing.T) *gin.Engine {
	router := gin.New()

	router.GET("/timeout", func(c *gin.Context) {
		ctx := c.Request.Context()
		select {
		case <-ctx.Done():
			fmt.Printf("setTimeout done\n")
		case <-time.After(2 * time.Second):
			assert.Fail(t, "test timeout fail")
		}
	})

	return router
}

func Test_Global_Timeout(t *testing.T) {
	router := setupDataFlow(t)

	const (
		longTimeout   = 400
		middleTimeout = 300
		shortTimeout  = 200
	)

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	// 只设置timeout
	SetTimeout(shortTimeout * time.Millisecond) //设置全局超时时间
	err := GET(ts.URL + "/timeout").Do()
	// 期望的结果是返回错误
	assert.Error(t, err)

	// 使用互斥api的原则，后面的覆盖前面的
	// 这里是WithContext生效, 预期超时时间400ms
	ctx, _ := context.WithTimeout(context.Background(), longTimeout*time.Millisecond)
	s := time.Now()
	SetTimeout(shortTimeout * time.Millisecond) // 设置全局超时时间
	err = GET(ts.URL + "/timeout").
		WithContext(ctx).
		Do()

	assert.Error(t, err)
	assert.GreaterOrEqual(t, int(time.Since(s)), int(middleTimeout*time.Millisecond))

	SetTimeout(time.Duration(0))
}

func Test_NewWithOpt_Timeout(t *testing.T) {
	router := setupDataFlow(t)

	const (
		longTimeout   = 400
		middleTimeout = 300
		shortTimeout  = 200
	)

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	// 只设置timeout
	greq := NewWithOpt(WithTimeout(shortTimeout * time.Millisecond)) //设置全局超时时间
	err := greq.GET(ts.URL + "/timeout").Do()
	// 期望的结果是返回错误
	assert.Error(t, err)

	// 使用互斥api的原则，后面的覆盖前面的
	// 这里是WithContext生效, 预期超时时间400ms
	ctx, _ := context.WithTimeout(context.Background(), longTimeout*time.Millisecond)
	s := time.Now()
	greq = NewWithOpt(WithTimeout(shortTimeout * time.Millisecond)) // 设置全局超时时间
	err = greq.GET(ts.URL + "/timeout").
		WithContext(ctx).
		Do()

	assert.Error(t, err)
	assert.GreaterOrEqual(t, int(time.Since(s)), int(middleTimeout*time.Millisecond))

}
