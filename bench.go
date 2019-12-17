package gout

import (
	"context"
	"github.com/guonaihong/gout/bench"
	"time"
)

// Bench provide benchmark features
type Bench struct {
	bench.Task

	df *DataFlow
}

// Concurrent set the number of benchmarks for concurrency
func (b *Bench) Concurrent(c int) *Bench {
	b.Task.Concurrent = c
	return b
}

// Number set the number of benchmarks
func (b *Bench) Number(n int) *Bench {
	b.Task.Number = n
	return b
}

// Rate set the frequency of the benchmark
func (b *Bench) Rate(rate int) *Bench {
	b.Task.Rate = rate
	return b
}

// Durations set the benchmark time
func (b *Bench) Durations(d time.Duration) *Bench {
	b.Task.Duration = d
	return b
}

// Do benchmark startup function
func (b *Bench) Do() error {
	// 报表插件
	req, err := b.df.Req.request()
	if err != nil {
		return err
	}

	client := b.df.out.Client
	if client == &DefaultClient {
		client = &DefaultBenchClient
	}

	r := bench.NewReport(context.Background(),
		b.Task.Concurrent,
		b.Task.Number,
		b.Task.Duration,
		req,
		client)

	// task是并发控制模块
	b.Run(r)
	return nil
}
