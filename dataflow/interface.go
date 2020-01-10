package dataflow

import (
	"io"
	"time"
)

type NewFilter interface {
	New(*DataFlow) interface{}
}

type Bencher interface {
	Concurrent(c int) Bencher
	Number(n int) Bencher
	Rate(rate int) Bencher
	Durations(d time.Duration) Bencher
	Do() error
}

type Retry interface {
	Attempt(attempt int) Retry
	WaitTime(waitTime time.Duration) Retry
	MaxWaitTime(maxWaitTime time.Duration) Retry
	Do() error
}

type Filter interface {
	Bench() Bencher
	Retry() Retry
}

type Curl interface {
	LongOption() Curl
	GenAndSend() Curl
	SetOutput(w io.Writer) Curl
	Do() error
}

type Export interface {
	Curl()
}
