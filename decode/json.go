package decode

import (
	"encoding/json"
	"io"
)

type JsonDecode struct {
	obj interface{}
}

func NewJsonDecode(obj interface{}) *JsonDecode {
	if obj == nil {
		return nil
	}
	return &JsonDecode{obj: obj}
}

func (j *JsonDecode) Decode(r io.Reader) error {
	decode := json.NewDecoder(r)
	return decode.Decode(j.obj)
}
