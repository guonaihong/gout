package gout

import (
	"context"
	"errors"
	"fmt"
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

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Retry is the core data structure of the retry function
// https://amazonaws-china.com/cn/blogs/architecture/exponential-backoff-and-jitter/
type Retry struct {
	df          *DataFlow
	attempt     int // Maximum number of attempts
	currAttempt int
	maxWaitTime time.Duration
	waitTime    time.Duration
}

// Attempt set the number of retries
func (r *Retry) Attempt(attempt int) *Retry {
	r.attempt = attempt
	return r
}

// WaitTime sets the basic wait time
func (r *Retry) WaitTime(waitTime time.Duration) *Retry {
	r.waitTime = waitTime
	return r
}

// MaxWaitTime Sets the maximum wait time
func (r *Retry) MaxWaitTime(maxWaitTime time.Duration) *Retry {
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

	req, err := r.df.request()
	if err != nil {
		return err
	}

	tk := time.NewTimer(r.maxWaitTime)
	for i := 0; i < r.attempt; i++ {

		resp, err := r.df.out.Client.Do(req)
		if err == nil {
			return r.df.bind(req, resp)
		}

		sleep := r.getSleep()

		if r.df.out.opt.Debug {
			fmt.Printf("filter:retry #current attempt:%d, wait time %v\n", r.currAttempt, sleep)
		}

		tk.Reset(sleep)
		ctx := r.df.getContext()
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
