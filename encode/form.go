package encode

import (
	"bytes"
	"fmt"
	"github.com/guonaihong/gout/core"
	"io/ioutil"
	"mime/multipart"
	"path/filepath"
	"reflect"
	"strconv"
	"unsafe"
)

var _ Adder = (*FormEncode)(nil)

type FormEncode struct {
	*multipart.Writer
}

func NewFormEncode(b *bytes.Buffer) *FormEncode {
	return &FormEncode{Writer: multipart.NewWriter(b)}
}

func stringToBytes(s string) []byte {
	sp := *(*[2]uintptr)(unsafe.Pointer(&s))
	bp := [3]uintptr{sp[0], sp[1], sp[1]}
	return *(*[]byte)(unsafe.Pointer(&bp))
}

func toBytes(v reflect.Value) (all []byte, err error) {
	if s, ok := v.Interface().(string); ok {
		all = stringToBytes(s)
	} else if b, ok := v.Interface().([]byte); ok {
		all = b
	} else {
		return nil, fmt.Errorf("unkown type toBytes:%T", v)
	}
	return all, nil
}

func (f *FormEncode) formFileWrite(key string, v reflect.Value, openFile bool) (err error) {
	var all []byte
	if openFile {
		var fileName string
		if s, ok := v.Interface().(string); ok {
			fileName = s
		} else if b, ok := v.Interface().([]byte); ok {
			fileName = core.BytesToString(b)
		} else {
			return fmt.Errorf("unkown type formFileWrite:%T, openFile:%t", v, openFile)
		}

		if all, err = ioutil.ReadFile(fileName); err != nil {
			return err
		}
	} else {
		if all, err = toBytes(v); err != nil {
			return err
		}
	}

	part, err := f.CreateFormFile(key, filepath.Base(key))
	if err != nil {
		return err
	}

	_, err = part.Write(all)
	return err
}

func (f *FormEncode) mapFormFile(key string, v reflect.Value, sf reflect.StructField) (next bool, err error) {
	var all []byte

	switch v.Interface().(type) {
	case core.FormFile:
		all, err = ioutil.ReadFile(string(v.Interface().(core.FormFile)))
		if err != nil {
			return false, err
		}

	case core.FormMem:
		all = []byte(v.Interface().(core.FormMem))
	default:
		return true, nil
	}

	part, err := f.CreateFormFile(key, filepath.Base(key))
	if err != nil {
		return false, err
	}

	_, err = part.Write(all)
	return false, err
}

func (f *FormEncode) Add(key string, v reflect.Value, sf reflect.StructField) (err error) {
	formFile := sf.Tag.Get("form-file")
	formMem := sf.Tag.Get("form-mem")
	b := false

	next, err := f.mapFormFile(key, v, sf)
	if err != nil {
		return err
	}

	if !next {
		return nil
	}

	if len(formFile) > 0 {
		if b, err = strconv.ParseBool(formFile); err != nil {
			return err
		}
		if !b {
			return nil
		}

		return f.formFileWrite(key, v, b)

	}

	if len(formMem) > 0 {
		if b, err = strconv.ParseBool(formMem); err != nil {
			return err
		}
		if !b {
			return nil
		}

		return f.formFileWrite(key, v, false)
	}

	return f.formFieldWrite(key, v)
}

func (f *FormEncode) formFieldWrite(key string, v reflect.Value) error {
	part, err := f.CreateFormField(key)
	var all []byte

	if all, err = toBytes(v); err != nil {
		return err
	}

	_, err = part.Write(all)
	return err
}

func (f *FormEncode) End() error {
	return f.Close()
}

func (f *FormEncode) Name() string {
	return "form"
}
