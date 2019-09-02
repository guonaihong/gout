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
	val = core.LoopElem(val)
	switch t := val.Kind(); t {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	case reflect.Float32, reflect.Float64:
	case reflect.String:
	default:
		if _, ok := val.Interface().([]byte); !ok {
			return fmt.Errorf("type(%T) %s:", b.obj, core.ErrUnkownType.Error())
		}
	}

	v := valToStr(val, emptyField)
	_, err := io.WriteString(w, v)
	return err
}
