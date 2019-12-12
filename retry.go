package gout

import (
	"math"
	"math/rand"
	"time"
)

var (
	RetryWaitTime    = 200 * time.Millisecond
	RetryMaxWaitTime = 10 * time.Second
	RetryAttempt     = 1
)

// https://amazonaws-china.com/cn/blogs/architecture/exponential-backoff-and-jitter/
type Retry struct {
	g           *routerGroup
	attempt     int // Maximum number of attempts
	currAttempt int
	maxWaitTime time.Duration
	waitTime    time.Duration
}

func (r *Retry) Attempt(attempt int) *Retry {
	r.attempt = attempt
	return r
}

func (r *Retry) WaitTime(waitTime time.Duration) *Retry {
	r.waitTime = waitTime
	return r
}

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

func (r *Retry) sleep() time.Duration {
	temp := r.waitTime * time.Duration(math.Exp2(float64(r.currAttempt)))
	temp = r.min(r.maxWaitTime, temp)
	return time.Duration(int(temp)/2 + rand.Intn(int(temp)/2))
}

func (r *Retry) Do() (err error) {
	defer r.reset()
	r.init()

	for i := 0; i < r.attempt; i++ {
		err = r.g.Do()
		if err == nil {
			return nil
		}

		sleep := r.sleep()

		time.Sleep(sleep)

		r.currAttempt++
	}

	return
}
