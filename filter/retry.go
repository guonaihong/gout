package filter

import (
	"context"
	"errors"
	"fmt"
	"github.com/guonaihong/gout/dataflow"
	"math"
	"math/rand"
	"time"
)

var (
	// RetryWaitTime retry basic wait time
	RetryWaitTime = 200 * time.Millisecond
	// RetryMaxWaitTime Maximum retry wait time
	RetryMaxWaitTime = 10 * time.Second
	// RetryAttempt number of retries
	RetryAttempt = 1
)

var (
	ErrRetryFail = errors.New("retry fail")
)

// Retry is the core data structure of the retry function
// https://amazonaws-china.com/cn/blogs/architecture/exponential-backoff-and-jitter/
type Retry struct {
	df          *dataflow.DataFlow
	attempt     int // Maximum number of attempts
	currAttempt int
	maxWaitTime time.Duration
	waitTime    time.Duration
}

func (r *Retry) New(df *dataflow.DataFlow) interface{} {
	return &Retry{df: df}
}

// Attempt set the number of retries
func (r *Retry) Attempt(attempt int) dataflow.Retry {
	r.attempt = attempt
	return r
}

// WaitTime sets the basic wait time
func (r *Retry) WaitTime(waitTime time.Duration) dataflow.Retry {
	r.waitTime = waitTime
	return r
}

// MaxWaitTime Sets the maximum wait time
func (r *Retry) MaxWaitTime(maxWaitTime time.Duration) dataflow.Retry {
	r.maxWaitTime = maxWaitTime
	return r
}

func (r *Retry) reset() {
	r.currAttempt = 0
}

func (r *Retry) init() {
	if r.attempt == 0 {
		r.attempt = RetryAttempt
	}

	if r.waitTime == 0 {
		r.waitTime = RetryWaitTime
	}

	if r.maxWaitTime == 0 {
		r.maxWaitTime = RetryMaxWaitTime
	}
}

// Does not pollute the namespace
func (r *Retry) min(a, b time.Duration) time.Duration {
	if a > b {
		return b
	}
	return a
}

func (r *Retry) getSleep() time.Duration {
	temp := r.waitTime * time.Duration(math.Exp2(float64(r.currAttempt)))
	if temp <= 0 {
		temp = r.waitTime
	}
	temp = r.min(r.maxWaitTime, temp)
	temp /= 2
	return temp + time.Duration(rand.Intn(int(temp)))
}

// Do send function
func (r *Retry) Do() (err error) {
	defer r.reset()
	r.init()

	req, err := r.df.Request()
	if err != nil {
		return err
	}

	tk := time.NewTimer(r.maxWaitTime)
	client := r.df.Client()

	for i := 0; i < r.attempt; i++ {

		// 这里不使用DataFlow.Do()方法原因是为了效率
		// 只需经过一次编码器得到request,后面就是多次使用
		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			return r.df.Bind(req, resp)
		}

		sleep := r.getSleep()

		if r.df.IsDebug() {
			fmt.Printf("filter:retry #current attempt:%d, wait time %v\n", r.currAttempt, sleep)
		}

		tk.Reset(sleep)
		ctx := r.df.GetContext()
		if ctx == nil {
			ctx = context.Background()
		}

		select {
		case <-tk.C:
			// 外部可以使用context直接取消
		case <-ctx.Done():
			return ctx.Err()
		}

		r.currAttempt++
	}

	return ErrRetryFail
}
