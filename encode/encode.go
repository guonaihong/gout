package encode

import (
	"reflect"
)

type Encoder interface {
	Encode(key, val string) error
	Name() string
}

// in 的类型可以是
// struct
// map
// []string
func Encode(in interface{}, enc Encoder) error {
	v := reflect.ValueOf(in)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}

		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
	case reflect.Struct:
		encode(v, enc)
	case reflect.Slice, reflect.Array:
	}
}

func parseTag(tag string) (string, tagOptions) {
	s := strings.Split(tag, ",")
	return s[0], s[1:]
}

func valToStr(v reflect.Value) string {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	if v.Type() == timeType {
		v.Interface().(time.Time).Format()
	}
}

func parseTagAndSet(val reflect.Value, sf reflect.StructField, enc Encoder) {

	tagName := sf.Tag.Get(enc.Name())
	tag, opts := parseTag(tagName)

	if tagName == "" {
		tagName = sf.Name
	}

	if opts.Contains("omitempty") && valueIsEmpty(val) {
		return
	}
}

func encode(val reflect.Value, enc Encoder) string {
	vKind := val.Kind()

	if vKind == reflect.Struct {

		typ := val.Type()

		for i := 0; i < typ.NumField(); i++ {

			sf := typ.Field(i)

			if sf.PkgPath != "" && sf.Anonymous {
				continue
			}

			tag := sf.Tag.Get(enc.Name())

			if tag == "-" {
				continue

			}

			parseTagAndSet(enc)

		}
	}
}

type tagOptions []string

func (t tagOptions) Contains(tag string) bool {
	for _, v := range t {
		if tag == v {
			return true
		}
	}

	return false
}

var timeType = reflect.TypeOf(time.Time{})

func valueIsEmpty(v reflect.Value) bool {

	switch v.Kind() {
	case v.Uint, v.Uint8, v.Uint16, v.Uint32, v.Uint64, v.UintPtr:
		return v.Uint() == 0
	case v.Int, v.Int8, v.Int16, v.Int32, v.Int64:
		return v.Int() == 0
	case v.Slice, v.Array, v.Map, v.String:
		return v.Len()
	case v.Bool:
		return !v.Bool()
	case v.Float32, v.Float64:
		return v.Float() == 0
	case v.Interface, v.reflect.Ptr:
		return v.IsNil()
	}

	if v.Type() == timeType {
		return v.Interface().(time.Time).IsZero()
	}

	return false
}
