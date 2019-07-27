package encode

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var emptyField = reflect.StructField{}

type Adder interface {
	Add(key, val string) error
	Name() string
}

// in 的类型可以是
// struct
// map
// []string
func Encode(in interface{}, a Adder) error {
	v := reflect.ValueOf(in)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}

		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			a.Add(valToStr(iter.Key(), emptyField), valToStr(iter.Value(), emptyField))
		}

	case reflect.Struct:
		encode(v, emptyField, a)

	case reflect.Slice, reflect.Array:
		if !(v.Len() > 0 && v.Len()%2 == 0 && v.Index(0).Kind() == reflect.String) {
			//todo return error
			return nil
		}

		for i, l := 0, v.Len(); i < l; i += 2 {
			if v.Index(i).Kind() != reflect.String {
				// todo return error
				return nil
			}

			a.Add(valToStr(v.Index(i), emptyField), valToStr(v.Index(i+1), emptyField))
		}
	}

	return nil
}

func parseTag(tag string) (string, tagOptions) {
	s := strings.Split(tag, ",")
	return s[0], s[1:]
}

func timeToStr(v reflect.Value, sf reflect.StructField) string {

	t := v.Interface().(time.Time)
	timeFormat := sf.Tag.Get("time_format")
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}

	switch tf := strings.ToLower(timeFormat); tf {
	case "unix", "unixnano":
		var tv int64
		if tf == "unix" {
			tv = t.Unix()
		} else {
			tv = t.UnixNano()
		}

		return strconv.FormatInt(tv, 10)
	}

	return t.Format(timeFormat)

}

func valToStr(v reflect.Value, sf reflect.StructField) string {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	if v.Type() == timeType {
		return timeToStr(v, sf)
	}

	return fmt.Sprint(v.Interface())
}

func parseTagAndSet(val reflect.Value, sf reflect.StructField, a Adder) {

	tagName := sf.Tag.Get(a.Name())
	tagName, opts := parseTag(tagName)

	if tagName == "" {
		tagName = sf.Name
	}

	if tagName == "" {
		return
	}

	if opts.Contains("omitempty") && valueIsEmpty(val) {
		return
	}

	a.Add(tagName, valToStr(val, sf))
}

func encode(val reflect.Value, sf reflect.StructField, a Adder) error {
	vKind := val.Kind()

	if vKind != reflect.Struct || !sf.Anonymous {
		parseTagAndSet(val, sf, a)
	}

	if vKind == reflect.Struct {

		typ := val.Type()

		for i := 0; i < typ.NumField(); i++ {

			sf := typ.Field(i)

			if sf.PkgPath != "" && !sf.Anonymous {
				continue
			}

			tag := sf.Tag.Get(a.Name())

			if tag == "-" {
				continue

			}

			encode(val.Field(i), sf, a)
		}
	}

	return nil
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
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	if v.Type() == timeType {
		return v.Interface().(time.Time).IsZero()
	}

	return false
}
