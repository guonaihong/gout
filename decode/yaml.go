package decode

import (
	"gopkg.in/yaml.v2"
	"io"
)

type YAMLDecode struct {
	obj interface{}
}

func NewYAMLDecode(obj interface{}) *YAMLDecode {
	if obj == nil {
		return nil
	}
	return &YAMLDecode{obj: obj}
}

func (x *YAMLDecode) Decode(r io.Reader) error {
	decode := yaml.NewDecoder(r)
	return decode.Decode(x.obj)
}

func DecodeYAML(r io.Reader, obj interface{}) error {
	decode := yaml.NewDecoder(r)
	return decode.Decode(obj)
}
