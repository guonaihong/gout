package gout

import (
	"context"
	"github.com/guonaihong/gout/bench"
	"time"
)

type Bench struct {
	bench.Task

	g *routerGroup
}

func (b *Bench) Concurrent(c int) *Bench {
	b.Task.Concurrent = c
	return b
}

func (b *Bench) Number(n int) *Bench {
	b.Task.Number = n
	return b
}

func (b *Bench) Rate(rate int) *Bench {
	b.Task.Rate = rate
	return b
}

func (b *Bench) Durations(d time.Duration) *Bench {
	b.Task.Duration = d
	return b
}

func (b *Bench) Do() error {
	// 报表插件
	r := bench.NewReport(context.Background(),
		b.Task.Concurrent,
		b.Task.Number,
		b.Task.Duration,
		"" /* todo */)

	// task是并发控制模块
	b.Run(r)
	return nil
}
