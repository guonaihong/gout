package encode

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

type formTest struct {
	set  interface{} //传入的值
	need interface{} //期望的值
	got  interface{}
}

func TestForm_FormEncodeNew(t *testing.T) {
	f := NewFormEncode(nil)
	assert.NotNil(t, f)
}

func TestForm_ToBytes(t *testing.T) {
	f := []formTest{
		{set: "test string", need: []byte("test string")},
		{set: []byte("test bytes"), need: []byte("test bytes")},
	}

	fail := []formTest{
		{set: time.Time{}},
		{set: time.Duration(0)},
	}

	// 测试正确的情况
	for _, v := range f {
		all, err := toBytes(reflect.ValueOf(v.set))
		assert.NoError(t, err)
		assert.Equal(t, all, v.need.([]byte))

	}

	// 测试错误的情况
	for _, v := range fail {
		_, err := toBytes(reflect.ValueOf(v))
		assert.Error(t, err)
	}
}

func TestForm_FormFileWrite(t *testing.T) {
}

func TestForm_MapFormFile(t *testing.T) {
}

func TestForm_Add(t *testing.T) {
}

func TestForm_FormFieldWrite(t *testing.T) {
}

func TestForm_End(t *testing.T) {
}
