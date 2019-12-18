package encode

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/textproto"
	"reflect"
	"strconv"
	"strings"

	"github.com/guonaihong/gout/core"
)

var _ Adder = (*FormEncode)(nil)

// FormEncode form-data encoder structure
type FormEncode struct {
	*multipart.Writer
}

// NewFormEncode create a new form-data encoder
func NewFormEncode(b *bytes.Buffer) *FormEncode {
	return &FormEncode{Writer: multipart.NewWriter(b)}
}

func toBytes(val reflect.Value) (all []byte, err error) {
	switch v := val.Interface().(type) {
	case string:
		all = core.StringToBytes(v)
	case []byte:
		all = v
	default:
		if val.Kind() == reflect.Interface {
			val = reflect.ValueOf(val.Interface())
		}

		switch t := val.Kind(); t {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		case reflect.Float32, reflect.Float64:
		case reflect.String:
		default:
			return nil, fmt.Errorf("unknown type toBytes:%T, kind:%v", v, val.Kind())
		}

		s := valToStr(val, emptyField)
		all = core.StringToBytes(s)
	}

	return all, nil
}

func (f *FormEncode) formFileWrite(key string, v reflect.Value, openFile bool) (err error) {
	var all []byte
	var contentType string
	var fileRealName = key
	if openFile {
		var fileName string
		switch v := v.Interface().(type) {
		case string:
			fileName = v
		case []byte:
			fileName = core.BytesToString(v)
		case core.FormType:
			s, ok := v.File.(string)
			if !ok {
				return fmt.Errorf("unknown type formFileWrite:%T, openFile:%t", v, openFile)
			}
			fileName = s
			fileRealName = v.FileName
			contentType = v.ContentType
		default:
			return fmt.Errorf("unknown type formFileWrite:%T, openFile:%t", v, openFile)
		}

		if all, err = ioutil.ReadFile(fileName); err != nil {
			return err
		}
	} else {
		switch val := v.Interface().(type) {
		case core.FormType:
			content, ok := val.File.([]byte)
			if !ok {
				s, okToo := val.File.(string)
				if !okToo {
					return fmt.Errorf("unknown type formFileWrite:%T, openFile:%t", v, openFile)
				}
				content = core.StringToBytes(s)
			}
			all = content
			contentType = val.ContentType
			fileRealName = val.FileName
		default:
			if all, err = toBytes(v); err != nil {
				return err
			}
		}

	}

	part, err := f.CreateFormFile(key, fileRealName, contentType)
	if err != nil {
		return err
	}

	_, err = part.Write(all)
	return err
}

func (f *FormEncode) mapFormFile(key string, v reflect.Value, sf reflect.StructField) (next bool, err error) {
	var all []byte
	var fileName = key
	var contentType string

	switch val := v.Interface().(type) {
	case core.FormType:

		fileName = val.FileName
		contentType = val.ContentType

		switch ft := val.File.(type) {
		case core.FormFile:
			all, err = ioutil.ReadFile(string(ft))
			if err != nil {
				return false, err
			}

		case core.FormMem:
			all = []byte(ft)
		default:
			return true, nil
		}

	case core.FormFile:
		all, err = ioutil.ReadFile(string(val))
		if err != nil {
			return false, err
		}

	case core.FormMem:
		all = []byte(val)
	default:
		return true, nil
	}

	part, err := f.CreateFormFile(key, fileName, contentType)
	if err != nil {
		return false, err
	}

	_, err = part.Write(all)
	return false, err
}

//下方为原函数附带的方法
var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// CreateFormFile 重写原来net/http里面的CreateFormFile函数
func (f *FormEncode) CreateFormFile(fieldName, fileName, contentType string) (io.Writer, error) {

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldName), escapeQuotes(fileName)))
	h.Set("Content-Type", contentType)
	return f.CreatePart(h)
}

// Add Encoder core function, used to set each key / value into the http form-data
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

// End refresh data
func (f *FormEncode) End() error {
	return f.Close()
}

// Name form-data Encoder name
func (f *FormEncode) Name() string {
	return "form"
}
