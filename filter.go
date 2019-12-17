package gout

// Filter 是过滤器核心函数
// 所有的过滤器都在此文件注册
// 过滤器要成为Filter的一个方法返回自己
type Filter struct {
	df *DataFlow
}

// Bench API performance stress test
func (f *Filter) Bench() *Bench {
	return &Bench{df: f.df}
}

// Retry API retry
func (f *Filter) Retry() *Retry {
	return &Retry{df: f.df}
}
