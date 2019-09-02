package core

import (
	"errors"
	"reflect"
)

type FormFile string

type FormMem []byte

var ErrUnkownType = errors.New("unkown type")

func LoopElem(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return v
		}
		v = v.Elem()
	}

	return v
}
