package encode

import (
	"bytes"
	"github.com/guonaihong/gout/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_WWWForm_New(t *testing.T) {
	assert.Nil(t, NewWWWFormEncode(nil))
	assert.NotNil(t, NewWWWFormEncode("h"))
}

type testWWWForm struct {
	w    *WWWFormEncode
	need string
}

func Test_WWWForm_Encode(t *testing.T) {
	var out bytes.Buffer

	tests := []testWWWForm{
		{NewWWWFormEncode(core.A{"k1", "v1", "k2", 2, "k3", 3.14}), "k1=v1&k2=2&k3=3.14"},
		{NewWWWFormEncode(core.H{"k1": "v1", "k2": 2, "k3": 3.14}), "k1=v1&k2=2&k3=3.14"},
	}

	for _, v := range tests {
		v.w.Encode(&out)
		assert.Equal(t, out.String(), v.need)
		out.Reset()
	}

	// 出错
	w := NewWWWFormEncode("fail")
	err := w.Encode(&out)
	assert.Error(t, err)
}

func Test_WWWForm_Name(t *testing.T) {
	assert.Equal(t, NewWWWFormEncode("").Name(), "www-form")
}
