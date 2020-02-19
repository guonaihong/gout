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

func Test_Form_genFormContext(t *testing.T) {
	fc := formContent{}
	f := []formTest{
		{set: "test string", need: []byte("test string")},
		{set: []byte("test bytes"), need: []byte("test bytes")},
		{set: interface{}(1), need: []byte("1")},
		{set: 1, need: []byte("1")},
	}

	fail := []formTest{
		{set: time.Time{}},
		//{set: time.Duration(0)},
	}

	// 测试正确的情况
	for _, v := range f {
		err := genFormContext("", reflect.ValueOf(v.set), emptyField, &fc)
		all := fc.data
		assert.NoError(t, err)
		assert.Equal(t, all, v.need.([]byte))

	}

	// 测试错误的情况
	for _, v := range fail {
		err := genFormContext("", reflect.ValueOf(v.set), emptyField, &fc)
		assert.Error(t, err)
	}
}

func checkForm(t *testing.T, boundary string, out *bytes.Buffer, caseID int) {
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
		if key == "voice" || key == "voice2" {
			assert.NotEqual(t, len(p.FileName()), 0, fmt.Sprintf("filename is empty:%d", caseID))
		}
		// slurp is value
		v := need[key]
		assert.Equal(t, v, string(slurp), fmt.Sprintf("fail test case id:%d", caseID))
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
	Voice  string `form:"voice" form-file:"xxx"`
	Voice2 string `form:"voice2" form-file:"true"`
}

//测试错误情况所用结构体2
type test_Form_struct_fail2 struct {
	Mode   string `form:"mode"`
	Text   string `form:"text"`
	Voice  string `form:"voice" form-file:"true"`
	Voice2 string `form:"voice2" form-file:"xxx"`
}

type test_Form_Second_struct_fail struct {
	Mode   string        `form:"mode"`
	Text   string        `form:"text"`
	Voice  core.FormType `form:"voice" form-file:"xxx"`
	Voice2 core.FormType `form:"voice2" form-file:"true"`
}

type test_Form_Second_struct_fail2 struct {
	Mode   string        `form:"mode"`
	Text   string        `form:"text"`
	Voice  core.FormType `form:"voice" form-mem:"true"`
	Voice2 core.FormType `form:"voice2" form-file:"xxx"`
}

type test_Form_Second_struct_Type_fail struct {
	Mode   string        `form:"mode"`
	Text   string        `form:"text"`
	Voice  core.FormType `form:"voice" form-file:"xxx"`
	Voice2 core.FormType `form:"voice2" form-file:"true"`
}

type test_Form_Second_struct_Inner_fail struct {
	Mode          string `form:"mode"`
	Text          string `form:"text"`
	core.FormType `form:"voice" form-file:"xxx"`
}

type test_Form_Second_struct_Inner_fail2 struct {
	Mode          string `form:"mode"`
	Text          string `form:"text"`
	core.FormType `form:"voice" form-file:"xxx"`
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
		{NewFormEncode(&out), core.H{
			"mode":   "A",
			"text":   "good",
			"voice":  core.FormMem("pcm1"),
			"voice2": core.FormFile("Non-existent file"), //不存在的文件
		}},
		{NewFormEncode(&out), core.H{
			"mode":   "A",
			"text":   "good",
			"voice":  core.FormType{FileName: "123.md", File: core.FormMem("pcm1")},
			"voice2": core.FormType{FileName: "123.md", File: core.FormFile("Non-existent file")}, //不存在的文件
		}},
		{NewFormEncode(&out), test_Form_struct_fail2{
			Mode:   "A",
			Text:   "good",
			Voice:  "pcm1",
			Voice2: "../testdata/voice.pcm",
		}},
		{NewFormEncode(&out), test_Form_struct_fail{ // id = 4
			Mode:   "A",
			Text:   "good",
			Voice:  "pcm1",
			Voice2: "../testdata/voice.pcm",
		}},
		{NewFormEncode(&out), test_Form_Second_struct_fail{ //id = 5
			Mode:   "A",
			Text:   "good",
			Voice:  core.FormType{FileName: "voice.pem", ContentType: "", File: "pcm1"},
			Voice2: core.FormType{FileName: "voice.pem", ContentType: "", File: "../testdata/voice.pcm"},
		}},
		{NewFormEncode(&out), test_Form_Second_struct_fail2{ //id = 6
			Mode:   "A",
			Text:   "good",
			Voice:  core.FormType{FileName: "voice.pem", ContentType: "", File: "pcm1"},
			Voice2: core.FormType{FileName: "voice.pem", ContentType: "", File: "../testdata/voice.pcm"},
		}},
		{NewFormEncode(&out), test_Form_Second_struct_Type_fail{ // id = 7
			Mode:   "A",
			Text:   "good",
			Voice:  core.FormType{FileName: "voice.pem", ContentType: "", File: 123},
			Voice2: core.FormType{FileName: "voice.pem", ContentType: "", File: "../testdata/voice.pcm"},
		}},
		{NewFormEncode(&out), test_Form_Second_struct_Type_fail{ // id = 8
			Mode:   "A",
			Text:   "good",
			Voice:  core.FormType{FileName: "voice.pem", ContentType: "", File: "pcm1"},
			Voice2: core.FormType{FileName: "voice.pem", ContentType: "", File: 123},
		}},
		{NewFormEncode(&out), test_Form_Second_struct_Inner_fail{
			"A",
			"good",
			core.FormType{FileName: "voice.pem", ContentType: "", File: "pcm1"}}},
		{NewFormEncode(&out), test_Form_Second_struct_Inner_fail2{
			"A",
			"good",
			core.FormType{FileName: "voice.pem", ContentType: "", File: "../testdata/voice.pcm"}}},
		{NewFormEncode(&out), test_Form_Third_struct2{ // id = 7
			"A",
			"good",
			core.FormType{FileName: "voice.pem", ContentType: "", File: []byte("Non-existent file")}}},
	}

	for id, v := range tests {
		err := Encode(v.data, v.f)
		assert.Error(t, err, fmt.Sprintf("test case id:%d", id))
	}
}

//测试正确情况所用结构体
type test_Form_struct struct {
	Mode   string `form:"mode"`
	Text   string `form:"text"`
	Voice  string `form:"voice" form-file:"mem"`
	Voice2 string `form:"voice2" form-file:"true"` // true 和file是相同的作用
}

//第二种测试情况
type test_Form_Second_struct struct {
	Mode   string        `form:"mode"`
	Text   string        `form:"text"`
	Voice  core.FormType `form:"voice" form-file:"mem"`
	Voice2 core.FormType `form:"voice2" form-file:"file"`
}

//第三种测试情况
type test_Form_Third_struct struct {
	Mode          string `form:"mode"`
	Text          string `form:"text"`
	core.FormType `form:"voice" form-file:"mem"`
}

type test_Form_Third_struct2 struct {
	Mode          string `form:"mode"`
	Text          string `form:"text"`
	core.FormType `form:"voice2" form-file:"true"`
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
			"voice2": core.FormType{FileName: "voice.pem", ContentType: "", File: core.FormFile("../testdata/voice.pcm")}}},

		{NewFormEncode(&out), test_Form_struct{
			Mode:   "A",
			Text:   "good",
			Voice:  "pcm1",
			Voice2: "../testdata/voice.pcm",
		}},
		{NewFormEncode(&out), test_Form_Second_struct{ // id = 3
			Mode:   "A",
			Text:   "good",
			Voice:  core.FormType{FileName: "voice.pem", ContentType: "", File: core.FormMem("pcm1")},
			Voice2: core.FormType{FileName: "voice.pem", ContentType: "", File: core.FormFile("../testdata/voice.pcm")},
		}},
		{NewFormEncode(&out), test_Form_Second_struct{ // id = 4
			Mode:   "A",
			Text:   "good",
			Voice:  core.FormType{FileName: "voice.pem", ContentType: "", File: "pcm1"},
			Voice2: core.FormType{FileName: "voice.pem", ContentType: "", File: "../testdata/voice.pcm"}}},
		{NewFormEncode(&out), test_Form_Third_struct{ // id = 5
			"A",
			"good",
			core.FormType{FileName: "voice.pem", ContentType: "", File: "pcm1"}}},
		{NewFormEncode(&out), test_Form_Third_struct2{ // id = 6
			"A",
			"good",
			core.FormType{FileName: "voice.pem", ContentType: "", File: "../testdata/voice.pcm"}}},
		{NewFormEncode(&out), test_Form_Third_struct2{ // id = 7
			"A",
			"good",
			core.FormType{FileName: "voice.pem", ContentType: "", File: []byte("../testdata/voice.pcm")}}},
	}

	for k, v := range tests {

		err := Encode(v.data, v.f)
		assert.NoError(t, err, fmt.Sprintf("test case id = %d", k))

		if err != nil {
			continue
		}
		v.f.End()

		boundary := v.f.Writer.Boundary()
		checkForm(t, boundary, &out, k)
		out.Reset()
	}
}

func Test_Form_Name(t *testing.T) {
	f := NewFormEncode(&bytes.Buffer{})
	assert.Equal(t, f.Name(), "form")
}
