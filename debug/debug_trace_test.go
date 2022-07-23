package debug

import (
	"bytes"
	"reflect"
	"strings"
)

func checkValue(b *bytes.Buffer) bool {
	info := &TraceInfo{}
	v := reflect.ValueOf(info)
	v = v.Elem()

	debugInfo := b.String()
	have := false
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		sf := typ.Field(i)
		if sf.PkgPath != "" {
			continue
		}

		name := sf.Name
		if !strings.Contains(debugInfo, name) {
			return false
		}
		have = true
	}

	if !have {
		return have
	}

	return true
}
