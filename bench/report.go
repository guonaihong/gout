package bench

import (
	"time"
)

// 数据字段，每个字段都用于显示
type report struct {
	Concurrency   int
	Failed        int
	Tps           float64
	Total         time.Duration
	Kbs           float64
	Mean          float64
	AllMean       float64
	Percentage55  time.Duration
	Percentage66  time.Duration
	Percentage75  time.Duration
	Percentage80  time.Duration
	Percentage90  time.Duration
	Percentage99  time.Duration
	Percentage100 time.Duration
}

type Report struct {
	report
}

func (r *Report) AddFail() {
	atomic.AddInt32(&r.Failed)
}
