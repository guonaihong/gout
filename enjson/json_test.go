package enjson

import (
	"bytes"
	"fmt"
	"github.com/guonaihong/gout/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testJSON struct {
	I int     `json:"i"`
	F float64 `json:"f"`
	S string  `json:"s"`
}

func TestMarshal(t *testing.T) {
	all, err := Marshal(map[string]interface{}{"a": "a"}, true)
	assert.NoError(t, err)
	assert.False(t, bytes.Contains(all, []byte("\n")))
}

func TestNewJSONEncode(t *testing.T) {
	j := NewJSONEncode(nil, false)
	assert.Nil(t, j)

	j = NewJSONEncode(nil, true)
	assert.Nil(t, j)
}

func TestJSONEncode_Name(t *testing.T) {
	assert.Equal(t, NewJSONEncode("", false).Name(), "json")
	assert.Equal(t, NewJSONEncode("", true).Name(), "json")
}

func TestJSONEncode_Encode(t *testing.T) {
	need := testJSON{
		I: 100,
		F: 3.14,
		S: "test encode json",
	}

	out := bytes.Buffer{}

	s := `{"I" : 100, "F" : 3.14, "S":"test encode json"}`
	data := []interface{}{need, &need, s, []byte(s)}
	for _, v := range data {
		for _, on := range []bool{true, false} {

			j := NewJSONEncode(v, on)
			out.Reset()

			assert.NoError(t, j.Encode(&out))

			got := testJSON{}

			err := json.Unmarshal(out.Bytes(), &got)
			assert.NoError(t, err)
			assert.Equal(t, got, need)
		}
	}

	// test fail
	for _, v := range []interface{}{func() {}, `{"query":"value"`} {
		for _, on := range []bool{true, false} {
			j := NewJSONEncode(v, on)
			out.Reset()
			err := j.Encode(&out)
			assert.Error(t, err, fmt.Sprintf("on:%t", on))
		}
	}
}
