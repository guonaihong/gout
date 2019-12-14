package encode

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"reflect"
	"testing"
	"time"

	"github.com/guonaihong/gout/core"
	"github.com/stretchr/testify/assert"
)

type formTest struct {
	set      interface{} //传入的值
	need     interface{} //期望的值
	got      interface{} //获取的值
	openFile bool
}

func Test_Form_FormEncodeNew(t *testing.T) {
	f := NewFormEncode(nil)
	assert.NotNil(t, f)
}

func Test_Form_toBytes(t *testing.T) {
	f := []formTest{
		{set: "test string", need: []byte("test string")},
		{set: []byte("test bytes"), need: []byte("test bytes")},
		{set: interface{}(1), need: []byte("1")},
		{set: 1, need: []byte("1")},
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

func Test_Form_FormFileWrite(t *testing.T) {

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

func checkForm(t *testing.T, boundary string, out *bytes.Buffer) {
	need := map[string]string{
		"mode":   "A",
		"text":   "good",
		"voice":  "pcm1",
		"voice2": "pcmpcmpcm\n",
	}

	mr := multipart.NewReader(out, boundary)

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}

		assert.NoError(t, err)

		slurp, err := ioutil.ReadAll(p)
		assert.NoError(t, err, fmt.Sprintf("formname = %s", p.FormName()))
		if err != nil {
			return
		}

		// key
		key := p.FormName()
		// slurp is value

		v := need[key]
		assert.Equal(t, v, string(slurp))
	}
}

type test_Form struct {
	f    *FormEncode
	data interface{}
}

//测试错误情况所用结构体
type test_Form_struct_fail struct {
	Mode   string `form:"mode"`
	Text   string `form:"text"`
	Voice  string `form:"voice" form-mem:"xxx"`
	Voice2 string `form:"voice2" form-file:"true"`
}

//测试错误情况所用结构体2
type test_Form_struct_fail2 struct {
	Mode   string `form:"mode"`
	Text   string `form:"text"`
	Voice  string `form:"voice" form-mem:"true"`
	Voice2 string `form:"voice2" form-file:"xxx"`
}

type test_Form_Second_struct_fail struct {
	Mode   string        `form:"mode"`
	Text   string        `form:"text"`
	Voice  core.FormType `form:"voice" form-mem:"xxx"`
	Voice2 core.FormType `form:"voice2" form-file:"true"`
}

type test_Form_Second_struct_fail2 struct {
	Mode   string        `form:"mode"`
	Text   string        `form:"text"`
	Voice  core.FormType `form:"voice" form-mem:"true"`
	Voice2 core.FormType `form:"voice2" form-file:"xxx"`
}
type test_Form_Second_struct_fail_Type struct {
	Mode   string        `form:"mode"`
	Text   string        `form:"text"`
	Voice  core.FormType `form:"voice" form-mem:"true"`
	Voice2 core.FormType `form:"voice2" form-file:"true"`
}

// 测试错误的情况
func Test_Form_Fail(t *testing.T) {
	var out bytes.Buffer
	tests := []test_Form{
		{NewFormEncode(&out), core.H{
			"mode":   "A",
			"text":   "good",
			"voice":  core.FormMem("pcm1"),
			"voice2": core.FormFile("Non-existent file"), //不存在的文件
		}},
		{NewFormEncode(&out), test_Form_struct_fail2{
			Mode:   "A",
			Text:   "good",
			Voice:  "pcm1",
			Voice2: "../testdata/voice.pcm",
		}},
		{NewFormEncode(&out), test_Form_struct_fail{
			Mode:   "A",
			Text:   "good",
			Voice:  "pcm1",
			Voice2: "../testdata/voice.pcm",
		}},
		{NewFormEncode(&out), test_Form_Second_struct_fail{
			Mode:   "A",
			Text:   "good",
			Voice:  core.FormType{FileName: "voice.pem", ContentType: "", File: "pcm1"},
			Voice2: core.FormType{FileName: "voice.pem", ContentType: "", File: "../testdata/voice.pcm"},
		}},
		{NewFormEncode(&out), test_Form_Second_struct_fail2{
			Mode:   "A",
			Text:   "good",
			Voice:  core.FormType{FileName: "voice.pem", ContentType: "", File: "pcm1"},
			Voice2: core.FormType{FileName: "voice.pem", ContentType: "", File: "../testdata/voice.pcm"},
		}},
		{NewFormEncode(&out), test_Form_Second_struct_fail_Type{
			Mode:   "A",
			Text:   "good",
			Voice:  core.FormType{FileName: "voice.pem", ContentType: "", File: 123},
			Voice2: core.FormType{FileName: "voice.pem", ContentType: "", File: "../testdata/voice.pcm"},
		}},
		{NewFormEncode(&out), test_Form_Second_struct_fail_Type{
			Mode:   "A",
			Text:   "good",
			Voice:  core.FormType{FileName: "voice.pem", ContentType: "", File: "pcm1"},
			Voice2: core.FormType{FileName: "voice.pem", ContentType: "", File: 123},
		}},
	}

	for _, v := range tests {
		err := Encode(v.data, v.f)
		assert.Error(t, err)
	}
}

//测试正确情况所用结构体
type test_Form_struct struct {
	Mode   string `form:"mode"`
	Text   string `form:"text"`
	Voice  string `form:"voice" form-mem:"true"`
	Voice2 string `form:"voice2" form-file:"true"`
}

//第二种测试情况
type test_Form_Second_struct struct {
	Mode   string        `form:"mode"`
	Text   string        `form:"text"`
	Voice  core.FormType `form:"voice" form-mem:"true"`
	Voice2 core.FormType `form:"voice2" form-file:"true"`
}

// 测试正确的情况
func Test_Form(t *testing.T) {
	var out bytes.Buffer

	tests := []test_Form{
		{NewFormEncode(&out), core.H{
			"mode":   "A",
			"text":   "good",
			"voice":  core.FormMem("pcm1"),
			"voice2": core.FormFile("../testdata/voice.pcm"),
		}},
		{NewFormEncode(&out), core.H{
			"mode":   "A",
			"text":   "good",
			"voice":  core.FormType{FileName: "voice.pem", ContentType: "", File: core.FormMem("pcm1")},
			"voice2": core.FormType{FileName: "voice.pem", ContentType: "", File: core.FormFile("../testdata/voice.pcm")},}},

		{NewFormEncode(&out), test_Form_struct{
			Mode:   "A",
			Text:   "good",
			Voice:  "pcm1",
			Voice2: "../testdata/voice.pcm",
		}},
		{NewFormEncode(&out), test_Form_Second_struct{
			Mode:   "A",
			Text:   "good",
			Voice:  core.FormType{FileName: "voice.pem", ContentType: "", File: core.FormMem("pcm1")},
			Voice2: core.FormType{FileName: "voice.pem", ContentType: "", File: core.FormFile("../testdata/voice.pcm")},},},
		{NewFormEncode(&out), test_Form_Second_struct{
			Mode:   "A",
			Text:   "good",
			Voice:  core.FormType{FileName: "voice.pem", ContentType: "", File: "pcm1"},
			Voice2: core.FormType{FileName: "voice.pem", ContentType: "", File: "../testdata/voice.pcm"},},},
	}

	for _, v := range tests {
		err := Encode(v.data, v.f)
		assert.NoError(t, err)

		if err != nil {
			continue
		}
		v.f.End()

		boundary := v.f.Writer.Boundary()
		checkForm(t, boundary, &out)
		out.Reset()
	}
}

func Test_Form_Name(t *testing.T) {
	f := NewFormEncode(&bytes.Buffer{})
	assert.Equal(t, f.Name(), "form")
}
