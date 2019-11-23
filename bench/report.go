package bench

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

var _ Tasker = (*Report)(nil)

type result struct {
	time       float64
	statusCode int
}

// 数据字段，每个字段都用于显示
type report struct {
	Concurrency   int   //并发数
	Failed        int32 //出错的连接数
	Tps           float64
	Duration      time.Duration // 连接总时间
	TotalBody     int           // 统计所有body大小
	TotalRead     int           // 统计所有read的流量
	SendNum       int           // 已经发送的http 请求
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
	StatusCodes   map[int]int
}

type Report struct {
	report
	Number    int // 发送总次数
	step      int // 动态报表输出间隔
	allResult chan result
	ctx       context.Context
	waitQuit  chan struct{} //等待startReport函数结束
	allTimes  []float64
}

func NewReport(ctx context.Context, c, n int, duration time.Duration, url string) *Report {
	step := 0
	if n > 150 {
		if step = n / 10; step < 100 {
			step = 10
		}
	}

	return &Report{
		allResult: make(chan result),
		report: report{
			Concurrency: c,
			StatusCodes: make(map[int]int, 2),
			Duration:    duration,
		},
		Number: n,
		step:   step,
		ctx:    ctx,
	}
}

// 初始化报表模块
func (r *Report) Init() {
	r.startReport()
}

// 负责构造压测http 链接和统计压测元数据
func (r *Report) SubProcess(work chan struct{}) {
}

// 等待结束
func (r *Report) WaitAll() {
	<-r.waitQuit
	r.outputReport() //输出最终报表
}

func (r *Report) addFail() {
	atomic.AddInt32(&r.Failed, 1)
}

func (r *Report) calBody(resp *http.Response) {

	/*
		bodyN := len(resp.Body)

		r.length = bodyN

		hN := len(resp.Status)
		hN += len(resp.Proto)
		hN += 1 //space
		hN += 2 //\r\n
		for k, v := range resp.Header {
			hN += len(k)

			for _, hv := range v {
				hN += len(hv)
			}
			hN += 2 //:space
			hN += 2 //\r\n
		}

		hN += 2

		atomic.AddInt32(&r.TotalBody, int32(bodyN))
		atomic.AddInt32(&r.TotalRead, int32(hN))
		atomic.AddInt32(&r.TotalRead, int32(bodyN))
	*/
}

func genTimeStr(now time.Time) string {
	year, month, day := now.Date()
	hour, min, sec := now.Clock()

	return fmt.Sprintf("%4d-%02d-%02d %02d:%02d:%02d.%06d",
		year,
		month,
		day,
		hour,
		min,
		sec,
		now.Nanosecond()/1e3,
	)
}

func (r *Report) startReport() {
	go func() {
		defer func() {
			fmt.Printf("  Finished  %15d requests\n", r.SendNum)
			r.waitQuit <- struct{}{}
		}()

		if r.step > 0 {
			for {
				select {
				case <-r.ctx.Done():
					return
				case v := <-r.allResult:
					r.SendNum++
					if r.step > 0 && r.SendNum%r.step == 0 {
						now := time.Now()

						fmt.Printf("    Opened %15d connections: [%s]\n",
							r.SendNum, genTimeStr(now))
					}

					r.allTimes = append(r.allTimes, v.time)
					r.StatusCodes[v.statusCode]++
				}
			}
		}

		begin := time.Now()
		interval := r.Duration / 10

		if interval == 0 || int(interval) > int(3*time.Second) {
			interval = 3 * time.Second
		}

		nTick := time.NewTicker(interval)
		count := 1
		for {
			select {
			case <-nTick.C:
				now := time.Now()

				fmt.Printf("  Completed %15d requests [%s]\n",
					r.SendNum, genTimeStr(now))

				count++
				next := begin.Add(time.Duration(count * int(interval)))
				if newInterval := next.Sub(time.Now()); newInterval > 0 {
					nTick = time.NewTicker(newInterval)
				} else {
					nTick = time.NewTicker(time.Millisecond * 100)
				}
			case v, ok := <-r.allResult:
				if !ok {
					return
				}

				r.SendNum++
				r.allTimes = append(r.allTimes, v.time)
				r.StatusCodes[v.statusCode]++
			case <-r.ctx.Done():
				return
			}
		}

	}()
}

func (r *Report) outputReport() {
}
