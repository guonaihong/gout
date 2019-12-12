package gout

//所有的过滤器都在此文件注册
type Filter struct {
	g *routerGroup
}

// API performance stress test
func (f *Filter) Bench() *Bench {
	return &Bench{g: f.g}
}

// API retry
func (f *Filter) Retry() *Retry {
	return &Retry{g: f.g}
}
