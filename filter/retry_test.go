package filter

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout/core"
	"github.com/guonaihong/gout/dataflow"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	retryCount        = 3
	retryDoesNotExist = ":6364"
)

func setupRetryFail() *gin.Engine {
	router := gin.New()

	var done chan struct{}
	router.GET("/", func(c *gin.Context) {
		<-done
	})

	return router
}

func setupRetryOk() *gin.Engine {
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
		assert.Equal(t, r.min(uint64(v.a), uint64(v.b)), uint64(v.need))
	}
}

func Test_Retry_sleep(t *testing.T) {
	r := Retry{attempt: 100, maxWaitTime: 10 * time.Second, waitTime: time.Second}
	r.init()

	// 方便画出曲线图
	for i := 0; i < r.attempt; i++ {
		cb := func() {
			sleep := r.getSleep()
			fmt.Printf("%d\n", sleep)
			//fmt.Printf("%d,%v\n", sleep, sleep)
			r.currAttempt++
		}
		b := assert.NotPanics(t, cb)
		if !b {
			break
		}
	}
}

func Test_Retry_init(t *testing.T) {
	r := Retry{attempt: RetryAttempt, maxWaitTime: RetryMaxWaitTime, waitTime: RetryWaitTime}
	r1 := Retry{}
	r1.init()
	assert.Equal(t, r1, r)
}

// test fail
func Test_Retry_fail(t *testing.T) {
	router := setupRetryOk()
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	tests := []dataflow.Retry{
		dataflow.POST(ts.URL + "/test.json").Debug(true).SetBody(&time.Time{}).Filter().Retry().Attempt(3),
	}

	for _, v := range tests {
		err := v.Do()
		assert.Error(t, err)
	}
}

func Test_Retry_Do(t *testing.T) {
	// 6364是随便写的一个端口，如果CI/CD这台机器上有这个端口，就需换个不存在的
	router := setupRetryFail()
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	urls := []string{ts.URL, retryDoesNotExist, retryDoesNotExist}
	// 测试全部超时的情况

	setTimeout := false
	for _, u := range urls {
		df := dataflow.GET(u)
		if u == retryDoesNotExist && !setTimeout || ts.URL == u {
			df.SetTimeout(11 * time.Millisecond)
			setTimeout = true
		}
		err := df.Debug(true).
			Filter().
			Retry().
			Attempt(retryCount).
			WaitTime(time.Millisecond * 10).
			MaxWaitTime(time.Millisecond * 50).
			Do()
		assert.Error(t, err)
	}

	// 测试正确的情况
	router = setupRetryOk()
	ts = httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	urls = []string{ts.URL}
	for _, u := range urls {
		err := dataflow.GET(u).
			// TODO 20 ms超时有时候会失败,分析下
			SetTimeout(30 * time.Millisecond).
			Debug(true).
			Filter().
			Retry().
			Attempt(retryCount).
			WaitTime(time.Millisecond * 10).
			MaxWaitTime(time.Millisecond * 50).
			Do()
		assert.NoError(t, err)
	}
}

func Test_Retry_Func(t *testing.T) {
	router := setupRetryOk()
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

	// test ok
	for _, err := range []error{
		// 演示使用备用端口
		func() error {
			//获取一个没有绑定服务的端口
			port := core.GetNoPortExists()
			s := ""

			err := dataflow.GET(":" + port).Debug(true).BindBody(&s).F().
				Retry().Attempt(3).WaitTime(time.Millisecond * 10).MaxWaitTime(time.Millisecond * 50).
				Func(func(c *dataflow.Context) error {
					if c.Error != nil {
						c.SetHost(ts.URL)
						return ErrRetry
					}
					return nil

				}).Do()

			assert.NoError(t, err)
			if err != nil {
				return err
			}
			assert.Equal(t, s, "ok")
			return nil
		}(),
		// 演示根据服务端不同错误码进行重试
		func() error {
			first := true
			// mock 服务端函数
			router := func() *gin.Engine {
				router := gin.New()

				router.GET("/code", func(c *gin.Context) {
					if first {
						c.String(209, "209")
						first = false
					} else {
						c.String(200, "ok")
					}
				})

				return router
			}()
			ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

			s := ""
			err := dataflow.GET(ts.URL + "/code").Debug(true).BindBody(&s).F().
				Retry().Attempt(3).WaitTime(time.Millisecond * 10).MaxWaitTime(time.Millisecond * 50).
				Func(func(c *dataflow.Context) error {
					if c.Error != nil || c.Code == 209 {
						return ErrRetry
					}

					return nil

				}).Do()

			assert.NoError(t, err)
			if err != nil {
				return err
			}
			assert.Equal(t, s, "ok")
			return nil
		}(),
	} {
		assert.NoError(t, err)
	}

	// test fail

	for _, err := range []error{
		func() error {

			s := ""
			err := dataflow.GET(ts.URL).Debug(true).BindBody(&s).F().
				Retry().Attempt(3).WaitTime(time.Millisecond * 10).MaxWaitTime(time.Millisecond * 50).
				Func(func(c *dataflow.Context) error {

					return errors.New("Do not retry")

				}).Do()
			assert.Error(t, err)
			if err != nil {
				return err
			}
			return nil
		}(),
		func() error {
			s := ""
			err := dataflow.GET(ts.URL).Debug(true).BindBody(&s).F().
				Retry().Attempt(3).WaitTime(time.Millisecond * 10).MaxWaitTime(time.Millisecond * 50).
				Func(func(c *dataflow.Context) error {
					// setbody不支持结构体，为了构造r.df.Request()返回错误
					c.SetBody(time.Time{})

					return ErrRetry

				}).Do()
			assert.Error(t, err)
			if err != nil {
				return err
			}
			return nil

		}(),
	} {
		assert.Error(t, err)
	}
}

// 测试retry返回error的情况
func Test_Filter_Retry_ReturnError(t *testing.T) {

	func() {
		var s string
		err := dataflow.New().POST(":" + core.GetNoPortExists()).BindJSON(&s).F().Retry().
			Attempt(-1).
			WaitTime(200 * time.Millisecond).
			MaxWaitTime(500 * time.Millisecond).Do()

		assert.Equal(t, err, ErrRetryFail)
		assert.True(t, errors.Is(err, ErrRetryFail))
	}()
	func() {
		var s string
		err := dataflow.New().POST(":" + core.GetNoPortExists()).BindJSON(&s).F().Retry().
			Attempt(3).
			WaitTime(200 * time.Millisecond).
			MaxWaitTime(500 * time.Millisecond).Do()

		assert.NotEqual(t, err, ErrRetryFail)
		assert.True(t, errors.Is(err, ErrRetryFail))
	}()
}

// 测试忽略io.EOF情况
func Test_Filter_Retry_ioEOF(t *testing.T) {
	router := func() *gin.Engine {
		router := gin.New()

		router.POST("/test/io/EOF", func(c *gin.Context) {
		})

		return router
	}()

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

	for id, err := range []error{
		func() error {
			var s string
			return dataflow.New().POST(":" + core.GetNoPortExists()).BindJSON(&s).F().Retry().
				Attempt(3).
				WaitTime(200 * time.Millisecond).
				MaxWaitTime(500 * time.Millisecond).
				Func(func(c *dataflow.Context) error {
					if c.Error != nil {
						c.SetURL(ts.URL + "/test/io/EOF")
						return ErrRetry
					}
					return nil
				}).Do()
		}(),
	} {
		assert.NoError(t, err, fmt.Sprintf("fail id:%d", id))
	}
}
