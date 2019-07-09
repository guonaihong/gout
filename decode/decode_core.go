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

// todo return value
func (d decData) Set(

	value reflect.Value,

	sf reflect.StructField,

	tagValue string) {

	setForm(d, value, sf, tagValue)
}

func setForm(m map[string][]string,
	value reflect.Value,
	sf reflect.StructField,
	tagValue string,
) {

	vs, ok := m[tagValue]
}

func decode(d decData, obj interface{}, tagName string) error {
	v := reflect.ValueOf(obj)
	if obj == nil || v.IsNil() {
		return errors.New("Wrong parameter")
	}

	decodeCore(d, reflect.StructField{}, v, tagName)
}

// todo delete
func parseTag(tag string) (string, tagOptions) {
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
	vKind := v.Kind()

	// elem pointer
	for vKind == reflect.Ptr {
		v = v.Elem()
	}

	if vKind == reflect.Struct || !sf.Anonymous {
		parseTagAndSet(val, sf, setter, tagName)
		return
	}

	if vKind == reflect.Struct {

		typ := val.Type()

		for i := 0; i < typ.NumField(); i++ {

			sf := typ.Field(i)

			if sf.PkgPath != "" && !sf.Anonymous {
				continue
			}

			tag := sf.Tag.Get()
			decodeCore(val.Field(i), sf, d, tag)
		}
	}

	return nil
}
