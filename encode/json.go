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
	//encode := json.NewEncoder(w)
	all, err := json.Marshal(j.obj)
	if err != nil {
		return err
	}

	// 不使用Encode函数的原因，encode结束之后会自作聪明的加'\n'
	//return encode.Encode(j.obj)
	_, err = w.Write(all)
	return err
}

// Name json Encoder name
func (j *JSONEncode) Name() string {
	return "json"
}
