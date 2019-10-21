package encode

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

type formTest struct {
	set      interface{} //传入的值
	need     interface{} //期望的值
	got      interface{} //获取的值
	openFile bool
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

	f := []formTest{
		{set: "../testdata/voice.pcm", need: nil, got: nil},
		{set: []byte("../testdata/voice.pcm"), need: nil, got: nil},
	}

	// TODO v0.0.3 换更好的测试策略，数据设置进去，再解析出来
	for _, v := range []bool{true, false} {
		for _, vv := range f {
			var out bytes.Buffer
			form := NewFormEncode(&out)
			assert.NotNil(t, form)

			err := form.formFileWrite("test form file write", reflect.ValueOf(vv.set), v)

			form.Close()
			assert.NoError(t, err)
			assert.NotEqual(t, out.Len(), 0)
		}
	}

	fail := []formTest{
		{set: "non-existent file", need: nil, got: nil, openFile: true},
		{set: time.Time{} /*不支持的类型*/, need: nil, got: nil, openFile: true},
		{set: time.Time{} /*不支持的类型*/, need: nil, got: nil, openFile: false},
	}

	for k, v := range fail {
		var out bytes.Buffer
		form := NewFormEncode(&out)
		assert.NotNil(t, form)

		err := form.formFileWrite("test form file write--fail", reflect.ValueOf(v.set), v.openFile)

		form.Close()
		assert.Error(t, err, fmt.Sprintf("index = %d", k))
		//assert.Equal(t, out.Len(), 0, fmt.Sprintf("index = %d:%s", k, out.Bytes()))

	}
}

func TestForm_MapFormFile(t *testing.T) {
}

func TestForm_Add(t *testing.T) {
}

func TestForm_FormFieldWrite(t *testing.T) {
}

func TestForm_End(t *testing.T) {
}
