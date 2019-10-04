package decode

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
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

func newValue(defValue interface{}) interface{} {
	p := reflect.New(reflect.TypeOf(defValue))
	p.Elem().Set(reflect.ValueOf(defValue))
	return p.Interface()
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
		{r: bytes.NewBufferString("1"), need: newValue(int(1)), got: new(int)},
		{r: bytes.NewBufferString("2"), need: newValue(int8(2)), got: new(int8)},
		{r: bytes.NewBufferString("3"), need: newValue(int16(3)), got: new(int16)},
		{r: bytes.NewBufferString("4"), need: newValue(int32(4)), got: new(int32)},
		{r: bytes.NewBufferString("5"), need: newValue(int64(5)), got: new(int64)},

		{r: bytes.NewBufferString("11"), need: newValue(uint(11)), got: new(uint)},
		{r: bytes.NewBufferString("12"), need: newValue(uint8(12)), got: new(uint8)},
		{r: bytes.NewBufferString("13"), need: newValue(uint16(13)), got: new(uint16)},
		{r: bytes.NewBufferString("14"), need: newValue(uint32(14)), got: new(uint32)},
		{r: bytes.NewBufferString("15"), need: newValue(uint64(15)), got: new(uint64)},

		{r: bytes.NewBufferString("3.14"), need: newValue(float32(3.14)), got: new(float32)},
		{r: bytes.NewBufferString("3.1415"), need: newValue(float64(3.1415)), got: new(float64)},

		{r: bytes.NewBufferString("test string"), need: newValue("test string"), got: new(string)},
		{r: bytes.NewBuffer([]byte("test bytes")), need: newValue([]byte("test bytes")), got: new([]byte)},
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
