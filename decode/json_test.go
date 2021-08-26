package decode

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewJSONDecode(t *testing.T) {
	j := NewJSONDecode(nil)
	assert.Nil(t, j)
}

type jsonTest struct {
	r    *bytes.Buffer
	need interface{}
	got  interface{}
}

func testDecodeJSON(t *testing.T, funcName string) {
	type jsonVal struct {
		A string `json:"a"`
		B string `json:"b"`
	}

	tests := []jsonTest{
		{r: bytes.NewBufferString(`{"a":"a", "b":"b"}`), need: &jsonVal{A: "a", B: "b"}, got: &jsonVal{}},
		{r: bytes.NewBufferString(`{"a":"aaa", "b":"bbb"}`), need: &jsonVal{A: "aaa", B: "bbb"}, got: &jsonVal{}},
	}

	for _, v := range tests {
		if funcName == "TestDecode" {
			j := NewJSONDecode(v.got)
			assert.NoError(t, j.Decode(v.r))
			assert.Equal(t, v.got, j.Value())
		} else {
			assert.NoError(t, JSON(v.r, v.got))
		}
		assert.Equal(t, v.need, v.got)
	}
}

func Test_json_DecodeJSON(t *testing.T) {
	testDecodeJSON(t, "TestDecodeJSON")
}

func Test_json_Decode(t *testing.T) {
	testDecodeJSON(t, "TestDecode")
}
