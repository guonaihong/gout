package encode

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testJSON struct {
	I int     `json:"i"`
	F float64 `json:"f"`
	S string  `json:"s"`
}

func TestNewJSONEncode(t *testing.T) {
	j := NewJSONEncode(nil)
	assert.Nil(t, j)

}

func TestJSONEncode_Name(t *testing.T) {
	assert.Equal(t, NewJSONEncode("").Name(), "json")
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
		j := NewJSONEncode(v)
		out.Reset()

		j.Encode(&out)

		got := testJSON{}

		err := json.Unmarshal(out.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, got, need)
	}

	// test fail
	for _, v := range []interface{}{func() {}, `{"query":"value"`} {
		j := NewJSONEncode(v)
		out.Reset()
		err := j.Encode(&out)
		assert.Error(t, err)
	}
}
