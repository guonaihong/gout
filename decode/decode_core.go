package decode

import (
	"net/http"
	"reflect"
)

type setter interface {
	Set(value reflect.Value,
		sf reflect.StructField,
		tagValue string) error
}

type decData map[string][]string

func (d decData) Set(
	value reflect.Value,
	sf reflect.StructField,
	tagValue string) {

}

func setForm(m map[string][]string, value reflect.Value, sf reflect.StructField, tagValue string) {
}

func decode(d decData, obj interface{}, tagName string) error {
	v := reflect.ValueOf(obj)
	if obj == nil || v.IsNil() {
		return errors.New("Wrong parameter")
	}

	decodeCore(d, reflect.StructField{}, v, tagName)
}

func decodeCore(val reflect.Value, sf reflect.StructField, d decData, tagName string) error {
	vKind := v.Kind()

	// elem pointer
	for vKind == reflect.Ptr {
		v = v.Elem()
	}

	if vKind == reflect.Struct || !sf.Anonymous {
		// todo set
	}

	if vKind == reflect.Struct {

		typ := val.Type()

		for i := 0; i < typ.NumField(); i++ {

			sf := typ.Field(i)

			if sf.PkgPath != "" && !sf.Anonymous {
				continue
			}

			tag := sf.Tag.Get()
			decodeCore(val.Field(i), sf, d)
		}
	}

	return nil
}
