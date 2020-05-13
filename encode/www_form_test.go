package encode

import (
	"bytes"
	"github.com/guonaihong/gout/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testWWWForm struct {
	w    *WWWFormEncode
	in   interface{}
	need string
}

func Test_WWWForm_Encode(t *testing.T) {
	var out bytes.Buffer

	tests := []testWWWForm{
		{NewWWWFormEncode(), core.A{"k1", "v1", "k2", 2, "k3", 3.14}, "k1=v1&k2=2&k3=3.14"},
		{NewWWWFormEncode(), core.H{"k1": "v1", "k2": 2, "k3": 3.14}, "k1=v1&k2=2&k3=3.14"},
	}

	for _, v := range tests {
		v.w.Encode(v.in)
		v.w.End(&out)
		assert.Equal(t, out.String(), v.need)
		out.Reset()
	}

}

func Test_WWWForm_Name(t *testing.T) {
	assert.Equal(t, NewWWWFormEncode().Name(), "www-form")
}
