package encode

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

type testYaml struct {
	I int     `yaml:"i"`
	F float64 `yaml:"f"`
	S string  `yaml:"s"`
}

func TestNewYAMLEncode(t *testing.T) {
	y := NewYAMLEncode(nil)
	assert.Nil(t, y)
}

func TestYAMLEncode_Name(t *testing.T) {
	assert.Equal(t, NewYAMLEncode("").Name(), "yaml")
}

func TestYAMLEncode_Encode(t *testing.T) {
	need := testYaml{
		I: 100,
		F: 3.14,
		S: "test encode yaml",
	}

	out := bytes.Buffer{}

	s := `
i: 100
f: 3.14
s: test encode yaml`
	data := []interface{}{need, &need, s, []byte(s)}
	for _, v := range data {
		x := NewYAMLEncode(v)
		out.Reset()

		x.Encode(&out)

		got := testYaml{}

		err := yaml.Unmarshal(out.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, got, need)
	}

	// test fail
	for _, v := range []interface{}{`I:100 {}`} {
		y := NewYAMLEncode(v)
		out.Reset()
		err := y.Encode(&out)
		assert.Error(t, err)
	}
}
