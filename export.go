package gout

// export 是特殊的过滤器
// 所有导出的过滤器都在次文件注册
type Export struct {
	df *DataFlow
}

func (e *Export) Curl() *Curl {
	return &Curl{}
}
