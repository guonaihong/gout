package encode

import (
	"reflect"
)

type Encoder interface {
	Encode(key, val string)
	Name() string
}

// in 的类型可以是
// struct
// map
// []string
func Encode(in interface{}, enc Encoder) error {
	v := reflect.ValueOf(in)

	switch v.Kind() {
	case reflect.Map:
	case reflect.Struct:
		encode(v, enc.Name())
	case reflect.Slice, reflect.Array:
	}
}

func parseTagAndSet() {
}

func encode(val reflect.Value, tagName string) string {
	vKind := val.Kind()

	if vKind == reflect.Struct {

		typ := val.Type()

		for i := 0; i < typ.NumField(); i++ {

			sf := typ.Field(i)

			if sf.PkgPath != "" && sf.Anonymous {
				continue
			}

			tag := sf.Tag.Get(tagName)

			if tag == "-" {
				continue

			}

			parseTagAndSet()

		}
	}
}
