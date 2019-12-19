package encode

import (
	"github.com/guonaihong/gout/core"
	"io"
	"net/url"
	"reflect"
)

var _ Adder = (*WWWFormEncode)(nil)

// WWWFormEncode x-www-form-urlencoded encoder structure
type WWWFormEncode struct {
	obj    interface{}
	values url.Values
}

// NewWWWFormEncode create a new x-www-form-urlencoded encoder
func NewWWWFormEncode(obj interface{}) *WWWFormEncode {
	if obj == nil {
		return nil
	}

	return &WWWFormEncode{obj: obj, values: make(url.Values)}
}

// Encode x-www-form-urlencoded encoder
func (we *WWWFormEncode) Encode(w io.Writer) (err error) {
	if err = Encode(we.obj, we); err != nil {
		return err
	}
	_, err = w.Write(core.StringToBytes(we.values.Encode()))
	return
}

// Add Encoder core function, used to set each key / value into the http x-www-form-urlencoded
// 这里value的设置暴露 reflect.Value和 reflect.StructField原因如下
// reflect.Value把value转成字符串
// reflect.StructField主要是可以在Add函数里面获取tag相关信息
func (we *WWWFormEncode) Add(key string, v reflect.Value, sf reflect.StructField) error {
	we.values.Add(key, valToStr(v, sf))
	return nil
}

// Name x-www-form-urlencoded Encoder name
func (we *WWWFormEncode) Name() string {
	return "www-form"
}
