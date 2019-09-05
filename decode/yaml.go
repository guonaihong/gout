package decode

import (
	"gopkg.in/yaml.v2"
	"io"
)

type YamlDecode struct {
	obj interface{}
}

func NewYamlDecode(obj interface{}) *YamlDecode {
	if obj == nil {
		return nil
	}
	return &YamlDecode{obj: obj}
}

func (x *YamlDecode) Decode(r io.Reader) error {
	decode := yaml.NewDecoder(r)
	return decode.Decode(x.obj)
}

func DecodeYAML(r io.Reader, obj interface{}) error {
	decode := yaml.NewDecoder(r)
	return decode.Decode(obj)
}
