package gout

import (
	"context"
	"github.com/guonaihong/gout/bench"
	"github.com/guonaihong/gout/dataflow"
	"time"
)

// Bench provide benchmark features
type Bench struct {
	bench.Task

	df *dataflow.DataFlow
}

// New
func (b *Bench) New(df *dataflow.DataFlow) interface{} {
	return &Bench{df: df}
}

// Concurrent set the number of benchmarks for concurrency
func (b *Bench) Concurrent(c int) dataflow.Bencher {
	b.Task.Concurrent = c
	return b
}

// Number set the number of benchmarks
func (b *Bench) Number(n int) dataflow.Bencher {
	b.Task.Number = n
	return b
}

// Rate set the frequency of the benchmark
func (b *Bench) Rate(rate int) dataflow.Bencher {
	b.Task.Rate = rate
	return b
}

// Durations set the benchmark time
func (b *Bench) Durations(d time.Duration) dataflow.Bencher {
	b.Task.Duration = d
	return b
}

// Do benchmark startup function
func (b *Bench) Do() error {
	// 报表插件
	req, err := b.df.Req.Request()
	if err != nil {
		return err
	}

	client := b.df.Client()
	if client == &dataflow.DefaultClient {
		client = &dataflow.DefaultBenchClient
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
