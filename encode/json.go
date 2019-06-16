package encode

import (
	"encoding/json"
	"io"
)

type JsonEncode struct {
	obj interface{}
}

func NewJsonEncode(obj interface{}) *JsonEncode {
	if obj == nil {
		return nil
	}

	return &JsonEncode{obj: obj}
}

func (j *JsonEncode) Encode(w io.Writer) error {
	encode := json.NewEncoder(w)
	return encode.Encode(j.obj)
}
