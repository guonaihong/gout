package encode

import (
	"encoding/json"
	"github.com/guonaihong/gout/core"
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
func (j *JSONEncode) Encode(w io.Writer) (err error) {
	if v, ok := core.GetBytes(j.obj); ok {
		if json.Valid(v) {
			_, err = w.Write(v)
			return err
		}
	}

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
