package decode

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewYAMLDecode(t *testing.T) {
	x := NewYAMLDecode(nil)
	assert.Nil(t, x)
}

type yamlTest struct {
	r    *bytes.Buffer
	need interface{}
	got  interface{}
}

func testDecodeYAML(t *testing.T, funcName string) {
	type yamlVal struct {
		A string `yaml:"a"`
		B string `yaml:"b"`
	}

	tests := []yamlTest{
		{r: bytes.NewBufferString(`a: a
b: b
`), need: &yamlVal{A: "a", B: "b"}, got: &yamlVal{}},
		{r: bytes.NewBufferString(`a: aaa
b: bbb
`), need: &yamlVal{A: "aaa", B: "bbb"}, got: &yamlVal{}},
	}

	for _, v := range tests {
		if funcName == "TestDecode" {
			x := NewYAMLDecode(v.got)
			x.Decode(v.r)
		} else {
			DecodeYAML(v.r, v.got)
		}
		assert.Equal(t, v.need, v.got)
	}
}

func Test_yaml_DecodeYAML(t *testing.T) {
	testDecodeYAML(t, "TestDecodeYAML")
}

func Test_yaml_Decode(t *testing.T) {
	testDecodeYAML(t, "TestDecode")
}
