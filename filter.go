package gout

//所有的过滤器都在此文件注册
type Filter struct {
	df *DataFlow
}

// API performance stress test
func (f *Filter) Bench() *Bench {
	return &Bench{df: f.df}
}

// API retry
func (f *Filter) Retry() *Retry {
	return &Retry{df: f.df}
}
