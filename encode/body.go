package encode

import (
	"fmt"
	"github.com/guonaihong/gout/core"
	"io"
	"reflect"
)

type BodyEncode struct {
	obj interface{}
}

func NewBodyEncode(obj interface{}) *BodyEncode {
	if obj == nil {
		return nil
	}

	return &BodyEncode{obj: obj}
}

func (b *BodyEncode) Encode(w io.Writer) error {
	val := reflect.ValueOf(b.obj)
	switch t := val.Kind(); t {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Interface:
		if _, ok := val.Interface().([]byte); !ok {
			return fmt.Errorf("type(%T) %s:", b.obj, core.ErrUnkownType.Error())
		}
	}

	v := valToStr(val, emptyField)
	_, err := io.WriteString(w, v)
	return err
}
