package encode

import (
	"gopkg.in/yaml.v2"
	"io"
)

type YamlEncode struct {
	obj interface{}
}

func NewYamlEncode(obj interface{}) *YamlEncode {
	if obj == nil {
		return nil
	}

	return &YamlEncode{obj: obj}
}

func (y *YamlEncode) Encode(w io.Writer) error {
	encode := yaml.NewEncoder(w)
	return encode.Encode(y.obj)
}
