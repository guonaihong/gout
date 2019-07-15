package decode

import (
	//"net/http"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type setter interface {
	Set(value reflect.Value,

		sf reflect.StructField,

		tagValue string) error
}

type decData map[string][]string

var emptyField = reflect.StructField{}

// todo return value
func (d decData) Set(

	value reflect.Value,

	sf reflect.StructField,

	tagValue string) error {

	return setForm(d, value, sf, tagValue)
}

func setForm(m map[string][]string,
	value reflect.Value,
	sf reflect.StructField,
	tagValue string,
) error {

	vs, ok := m[tagValue]
	if !ok {
		return nil
	}

	switch value.Kind() {
	case reflect.Slice:
		return setSlice(vs, sf, value)
	case reflect.Array:
		if len(vs) != value.Len() {
			return fmt.Errorf("%q is not valid value for %s", vs, value.Type().String())
		}

		return setArray(vs, sf, value)
	}

	var val string
	if len(vs) > 0 {
		val = vs[0]
	}

	return setBase(val, sf, value)
}

func decode(d decData, obj interface{}, tagName string) error {
	v := reflect.ValueOf(obj)
	if obj == nil || v.IsNil() {
		return errors.New("Wrong parameter")
	}

	return decodeCore(v, emptyField, d, tagName)
}

// todo delete
func parseTag(tag string) (string, []string) {
	s := strings.Split(tag, ",")
	return s[0], s[1:]
}

func parseTagAndSet(val reflect.Value, sf reflect.StructField, setter setter, tagName string) {
	tagName = sf.Tag.Get(tagName)
	tagName, _ = parseTag(tagName)

	if tagName == "" {
		tagName = sf.Name
	}

	if tagName == "" {
		return
	}

	setter.Set(val, sf, tagName)
}

func decodeCore(val reflect.Value, sf reflect.StructField, setter setter, tagName string) error {
	vKind := val.Kind()

	// elem pointer
	for vKind == reflect.Ptr {
		val = val.Elem()
	}

	if vKind == reflect.Struct || !sf.Anonymous {
		parseTagAndSet(val, sf, setter, tagName)
		return nil
	}

	if vKind == reflect.Struct {

		typ := val.Type()

		for i := 0; i < typ.NumField(); i++ {

			sf := typ.Field(i)

			if sf.PkgPath != "" && !sf.Anonymous {
				continue
			}

			tag := sf.Tag.Get(tagName)
			decodeCore(val.Field(i), sf, setter, tag)
		}
	}

	return nil
}

type convert struct {
	bitSize int
	cb      func(val string, bitSize int, sf reflect.StructField, field reflect.Value) error
}

var convertFunc = map[reflect.Kind]convert{
	reflect.Uint:   {bitSize: 0, cb: setIntField},
	reflect.Uint8:  {bitSize: 8, cb: setIntField},
	reflect.Uint16: {bitSize: 16, cb: setIntField},
	reflect.Uint32: {bitSize: 32, cb: setIntField},
	reflect.Uint64: {bitSize: 64, cb: setIntField},
	reflect.Int:    {bitSize: 0, cb: setIntField},
	reflect.Int8:   {bitSize: 8, cb: setUintField},
	reflect.Int16:  {bitSize: 16, cb: setUintField},
	reflect.Int32:  {bitSize: 32, cb: setUintField},
	reflect.Int64:  {bitSize: 64, cb: setUintField},
	reflect.Bool:   {bitSize: 0, cb: setBoolField},
}

func setIntField(val string, bitSize int, sf reflect.StructField, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, sf reflect.StructField, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(val string, bitSize int, sf reflect.StructField, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(val string, bitSize int, sf reflect.StructField, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setTimeField(val string, bitSize int, structField reflect.StructField, value reflect.Value) error {
	timeFormat := structField.Tag.Get("time_format")
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}

	switch tf := strings.ToLower(timeFormat); tf {
	case "unix", "unixnano":
		tv, err := strconv.ParseInt(val, 10, 0)
		if err != nil {
			return err
		}

		d := time.Duration(1)
		if tf == "unixnano" {
			d = time.Second
		}

		t := time.Unix(tv/int64(d), tv%int64(d))
		value.Set(reflect.ValueOf(t))
		return nil

	}

	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != "" {
		loc, err := time.LoadLocation(locTag)
		if err != nil {
			return err
		}
		l = loc
	}

	t, err := time.ParseInLocation(timeFormat, val, l)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(t))
	return nil
}

func setArray(vals []string, sf reflect.StructField, value reflect.Value) error {
	for i, s := range vals {
		err := setBase(s, sf, value.Index(i))
		if err != nil {
			return err
		}
	}
	return nil
}

func setSlice(vals []string, sf reflect.StructField, value reflect.Value) error {
	slice := reflect.MakeSlice(value.Type(), len(vals), len(vals))
	err := setArray(vals, sf, slice)
	if err != nil {
		return err
	}
	value.Set(slice)
	return nil
}

func setTimeDuration(val string, bitSize int, value reflect.Value) error {
	d, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(d))
	return nil
}

func setBase(val string, sf reflect.StructField, value reflect.Value) error {
	fn, ok := convertFunc[value.Kind()]
	if ok {
		fn.cb(val, fn.bitSize, sf, value)
		return nil
	}

	return nil
}
