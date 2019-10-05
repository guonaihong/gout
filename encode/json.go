package encode

import (
	"encoding/json"
	"io"
)

type JSONEncode struct {
	obj interface{}
}

func NewJSONEncode(obj interface{}) *JSONEncode {
	if obj == nil {
		return nil
	}

	return &JSONEncode{obj: obj}
}

func (j *JSONEncode) Encode(w io.Writer) error {
	encode := json.NewEncoder(w)
	return encode.Encode(j.obj)
}
