package encode

import (
	"gopkg.in/yaml.v2"
	"io"
)

type YAMLEncode struct {
	obj interface{}
}

func NewYAMLEncode(obj interface{}) *YAMLEncode {
	if obj == nil {
		return nil
	}

	return &YAMLEncode{obj: obj}
}

func (y *YAMLEncode) Encode(w io.Writer) error {
	encode := yaml.NewEncoder(w)
	return encode.Encode(y.obj)
}
