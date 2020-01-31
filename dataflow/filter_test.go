package dataflow

import (
	"github.com/guonaihong/gout/bench"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testBench struct{}

func (t *testBench) New(df *DataFlow) interface{} {
	return &testBench{}
}

func (t *testBench) Concurrent(c int) Bencher {
	return t
}

func (t *testBench) Number(n int) Bencher {
	return t
}

func (t *testBench) Rate(rate int) Bencher {
	return t
}
func (t *testBench) Durations(d time.Duration) Bencher {
	return t
}

func (t *testBench) Loop(func(c *Context) error) Bencher {
	return t
}

func (t *testBench) GetReport(r *bench.Report) Bencher {
	return t
}

func (t *testBench) Do() error {
	return nil
}

type testRetry struct{}

func (t *testRetry) New(df *DataFlow) interface{} {
	return &testRetry{}
}

func (t *testRetry) Attempt(attempt int) Retry {
	return t
}

func (t *testRetry) WaitTime(waitTime time.Duration) Retry {
	return t
}

func (t *testRetry) MaxWaitTime(maxWaitTime time.Duration) Retry {
	return t
}

func (t *testRetry) Func(func(c *Context) error) Retry {
	return t
}

func (t *testRetry) Do() error {
	return nil
}

const (
	benchName = "bench"
	retryName = "retry"
)

func Test_Filter_Bench(t *testing.T) {
	bkcurl, ok := filters[benchName]
	delete(filters, benchName)
	defer func() {
		if ok {
			filters[benchName] = bkcurl
		}
	}()

	// test panic
	for _, v := range []func(){
		func() {
			f := filter{}
			f.Bench()
		},
		func() {
			filters[benchName] = &testCurlFail{}
			f := filter{}
			f.Bench()
		},
	} {
		assert.Panics(t, v)
	}

	//test ok
	for _, v := range []func(){
		func() {
			filters[benchName] = &testBench{}
			f := filter{}
			f.Bench()
		},
	} {
		assert.NotPanics(t, v)
	}
}

func Test_Filter_Retry(t *testing.T) {
	bkcurl, ok := filters[retryName]
	delete(filters, retryName)
	defer func() {
		if ok {
			filters[retryName] = bkcurl
		}
	}()

	// test panic
	for _, v := range []func(){
		func() {
			f := filter{}
			f.Retry()
		},
		func() {
			filters[retryName] = &testCurlFail{}
			f := filter{}
			f.Retry()
		},
	} {
		assert.Panics(t, v)
	}

	//test ok
	for _, v := range []func(){
		func() {
			filters[retryName] = &testRetry{}
			f := filter{}
			f.Retry()
		},
	} {
		assert.NotPanics(t, v)
	}
}
