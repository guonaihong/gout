package decode

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
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
			assert.NoError(t, x.Decode(v.r))
			assert.Equal(t, v.got, x.Value())
		} else {
			assert.NoError(t, YAML(v.r, v.got))
		}
		assert.Equal(t, v.need, v.got)
	}
}

func Test_yaml_YAML(t *testing.T) {
	testDecodeYAML(t, "TestYAML")
}

func Test_yaml_Decode(t *testing.T) {
	testDecodeYAML(t, "TestDecode")
}
