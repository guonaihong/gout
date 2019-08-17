package encode

import (
	"bytes"
	"fmt"
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

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b[0]))
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

func (f *FormEncode) partWrite(key string, v reflect.Value, openFile bool) (err error) {
	var all []byte
	if openFile {
		var fileName string
		if s, ok := v.Interface().(string); ok {
			fileName = s
		} else if b, ok := v.Interface().([]byte); ok {
			fileName = bytesToString(b)
		} else {
			return fmt.Errorf("unkown type partWrite:%T, openFile:%t", v, openFile)
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

func (f *FormEncode) Add(key string, v reflect.Value, sf reflect.StructField) (err error) {
	formFile := sf.Tag.Get("form-file")
	formMem := sf.Tag.Get("form-mem")
	b := false
	all := []byte{}

	if len(formFile) > 0 {
		if b, err = strconv.ParseBool(formFile); err != nil {
			return err
		}
		if !b {
			return nil
		}

		return f.partWrite(key, v, b)

	}

	if len(formMem) > 0 {
		if b, err = strconv.ParseBool(formMem); err != nil {
			return err
		}
		if !b {
			return nil
		}

		return f.partWrite(key, v, false)
	}

	part, err := f.CreateFormField(key)

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
