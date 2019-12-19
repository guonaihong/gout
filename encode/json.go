package encode

import (
	"encoding/json"
	"io"
)

// JSONEncode json encoder structure
type JSONEncode struct {
	obj interface{}
}

// NewJSONEncode create a new json encoder
func NewJSONEncode(obj interface{}) *JSONEncode {
	if obj == nil {
		return nil
	}

	return &JSONEncode{obj: obj}
}

// Encode json encoder
func (j *JSONEncode) Encode(w io.Writer) error {
	encode := json.NewEncoder(w)
	return encode.Encode(j.obj)
}

// Name json Encoder name
func (j *JSONEncode) Name() string {
	return "json"
}
