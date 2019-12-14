package gout

//所有的过滤器都在此文件注册
type Filter struct {
	df *DataFlow
}

func (f *Filter) Bench() *Bench {
	return &Bench{df: f.df}
}
