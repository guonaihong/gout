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

	switch t := reflect.ValueOf(b.obj).Kind(); t {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Interface:
		return fmt.Errorf("type(%T) %s:", t, core.ErrUnkownType.Error())
	}

	v := valToStr(reflect.ValueOf(b.obj), emptyField)
	_, err := io.WriteString(w, v)
	return err
}
