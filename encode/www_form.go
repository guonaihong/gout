package encode

import (
	"github.com/guonaihong/gout/core"
	"io"
	"net/url"
	"reflect"
)

var _ Adder = (*WWWFormEncode)(nil)

type WWWFormEncode struct {
	obj    interface{}
	values url.Values
}

func NewWWWFormEncode(obj interface{}) *WWWFormEncode {
	if obj == nil {
		return nil
	}

	return &WWWFormEncode{obj: obj, values: make(url.Values)}
}

func (we *WWWFormEncode) Encode(w io.Writer) (err error) {
	if err = Encode(we.obj, we); err != nil {
		return err
	}
	_, err = w.Write(core.StringToBytes(we.values.Encode()))
	return
}

func (we *WWWFormEncode) Add(key string, v reflect.Value, sf reflect.StructField) error {
	we.values.Add(key, valToStr(v, sf))
	return nil
}

func (we *WWWFormEncode) Name() string {
	return "www-form"
}
