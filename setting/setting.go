package setting

import "time"

// 设置
type Setting struct {
	// 控制是否使用空值
	NotIgnoreEmpty bool

	//是否自动加ContentType
	NoAutoContentType bool

	//超时时间
	Timeout time.Duration

	// 目前用作SetTimeout 和 WithContext是互斥
	// index是自增id，主要给互斥API定优先级
	// 对于互斥api，后面的会覆盖前面的
	Index int

	//当前time 的index
	TimeoutIndex int

	UseChunked bool
}

// 使用chunked数据
func (s *Setting) Chunked() {
	s.UseChunked = true
}

func (s *Setting) SetTimeout(d time.Duration) {
	s.Index++
	s.TimeoutIndex = s.Index
	s.Timeout = d
}

func (s *Setting) Reset() {
	s.NotIgnoreEmpty = false
	s.NoAutoContentType = false
	s.Index = 0
	//s.TimeoutIndex = 0
	s.Timeout = time.Duration(0)
}
