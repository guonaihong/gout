package decode

import (
	"bytes"
	"errors"
	"github.com/guonaihong/gout/core"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewBodyDecode(t *testing.T) {
	b := NewBodyDecode(nil)
	assert.Nil(t, b)

	b = NewBodyDecode(new(int))
	assert.NotNil(t, b)
}

type bodyTest struct {
	r    *bytes.Buffer
	need interface{}
	got  interface{}
}

type myFailRead struct{}

func (m *myFailRead) Read(p []byte) (n int, err error) {
	return 0, errors.New("fail")

}
func TestDecodeBodyFail(t *testing.T) {
	var tm time.Time
	err := DecodeBody(bytes.NewBufferString("1"), &tm)
	assert.Error(t, err)

	err = DecodeBody(&myFailRead{}, new(int))
	assert.Error(t, err)
}

func TestDecode(t *testing.T) {
	testDecodeBody(t, "TestDecode")
}

func TestDecodeBody(t *testing.T) {
	testDecodeBody(t, "TestDecodeBody")
}

func testDecodeBody(t *testing.T, funcName string) {
	tests := []bodyTest{
		{r: bytes.NewBufferString("1"), need: core.NewPtrVal(int(1)), got: new(int)},
		{r: bytes.NewBufferString("2"), need: core.NewPtrVal(int8(2)), got: new(int8)},
		{r: bytes.NewBufferString("3"), need: core.NewPtrVal(int16(3)), got: new(int16)},
		{r: bytes.NewBufferString("4"), need: core.NewPtrVal(int32(4)), got: new(int32)},
		{r: bytes.NewBufferString("5"), need: core.NewPtrVal(int64(5)), got: new(int64)},

		{r: bytes.NewBufferString("11"), need: core.NewPtrVal(uint(11)), got: new(uint)},
		{r: bytes.NewBufferString("12"), need: core.NewPtrVal(uint8(12)), got: new(uint8)},
		{r: bytes.NewBufferString("13"), need: core.NewPtrVal(uint16(13)), got: new(uint16)},
		{r: bytes.NewBufferString("14"), need: core.NewPtrVal(uint32(14)), got: new(uint32)},
		{r: bytes.NewBufferString("15"), need: core.NewPtrVal(uint64(15)), got: new(uint64)},

		{r: bytes.NewBufferString("3.14"), need: core.NewPtrVal(float32(3.14)), got: new(float32)},
		{r: bytes.NewBufferString("3.1415"), need: core.NewPtrVal(float64(3.1415)), got: new(float64)},

		{r: bytes.NewBufferString("test string"), need: core.NewPtrVal("test string"), got: new(string)},
		{r: bytes.NewBuffer([]byte("test bytes")), need: core.NewPtrVal([]byte("test bytes")), got: new([]byte)},
	}

	for _, v := range tests {
		if funcName == "TestDecode" {
			body := NewBodyDecode(v.got)
			body.Decode(v.r)
		} else {
			DecodeBody(v.r, v.got)
		}
		assert.Equal(t, v.need, v.got)
	}
}
