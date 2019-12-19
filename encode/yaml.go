package encode

import (
	"gopkg.in/yaml.v2"
	"io"
)

// YAMLEncode yaml encoder structure
type YAMLEncode struct {
	obj interface{}
}

// NewYAMLEncode create a new yaml encoder
func NewYAMLEncode(obj interface{}) *YAMLEncode {
	if obj == nil {
		return nil
	}

	return &YAMLEncode{obj: obj}
}

// Encode yaml encoder
func (y *YAMLEncode) Encode(w io.Writer) error {
	encode := yaml.NewEncoder(w)
	return encode.Encode(y.obj)
}

// Name yaml Encoder name
func (y *YAMLEncode) Name() string {
	return "yaml"
}
