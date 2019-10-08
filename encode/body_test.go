package encode

import (
	"bytes"
	"github.com/guonaihong/gout/core"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestNewBodyEncode(t *testing.T) {
	b := NewBodyEncode(nil)
	assert.Nil(t, b)
}

type bodyTest struct {
	w    *strings.Builder
	set  interface{}
	need string
}

func Test_body_EncodeFail(t *testing.T) {
	b := NewBodyEncode(&time.Time{})
	err := b.Encode(bytes.NewBufferString("hello world"))
	assert.Error(t, err)
}

func Test_body_Encode(t *testing.T) {

	tests := []bodyTest{
		{w: &strings.Builder{}, set: int(1), need: "1"},
		{w: &strings.Builder{}, set: int8(2), need: "2"},
		{w: &strings.Builder{}, set: int16(3), need: "3"},
		{w: &strings.Builder{}, set: int32(4), need: "4"},
		{w: &strings.Builder{}, set: int64(5), need: "5"},

		{w: &strings.Builder{}, set: uint(11), need: "11"},
		{w: &strings.Builder{}, set: uint8(12), need: "12"},
		{w: &strings.Builder{}, set: uint16(13), need: "13"},
		{w: &strings.Builder{}, set: uint32(14), need: "14"},
		{w: &strings.Builder{}, set: uint64(15), need: "15"},
		{w: &strings.Builder{}, set: []byte("test bytes"), need: "test bytes"},
		{w: &strings.Builder{}, set: "test string", need: "test string"},
		{w: &strings.Builder{}, set: int(1), need: "1"},
		{w: &strings.Builder{}, set: core.NewPtrVal(1010), need: "1010"},

		// test io.Reader
		{w: &strings.Builder{}, set: bytes.NewBufferString("set body:hello world"), need: "set body:hello world"},
	}

	for _, v := range tests {
		b := NewBodyEncode(v.set)
		b.Encode(v.w)
		assert.Equal(t, v.w.String(), v.need)
	}
}
