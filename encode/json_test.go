package encode

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testJson struct {
	I int     `json:"i"`
	F float64 `json:"f"`
	S string  `json:"s"`
}

func TestNewJSONEncode(t *testing.T) {
	j := NewJSONEncode(nil)
	assert.Nil(t, j)

}

func TestJSONEncode_Encode(t *testing.T) {
	need := testJson{
		I: 100,
		F: 3.14,
		S: "test encode json",
	}

	out := bytes.Buffer{}

	data := []interface{}{need, &need}
	for _, v := range data {
		j := NewJSONEncode(v)
		out.Reset()

		j.Encode(&out)

		got := testJson{}

		err := json.Unmarshal(out.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, got, need)
	}

}
