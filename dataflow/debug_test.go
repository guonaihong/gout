package dataflow

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/guonaihong/gout/color"
	"github.com/guonaihong/gout/core"
	"github.com/guonaihong/gout/debug"
	"github.com/stretchr/testify/assert"
)

// 测试Debug接口在各个数据类型下面是否可以打印request和response消息
func TestDebug_Debug(t *testing.T) {

	rspVal := "hello world"
	ts := createGeneral(rspVal)
	defer ts.Close()

	var buf bytes.Buffer

	dbug := func() debug.Apply {
		return debug.Func(func(o *debug.Options) {
			o.Debug = true
			o.Color = true
			o.Write = &buf
		})
	}()

	for index, err := range []error{
		// 测试http header里面带%的情况
		func() error {
			buf.Reset()
			err := New().POST(ts.URL).SetHeader(core.H{`h%`: "hello%", "Cookie": `username=admin; token=b7ea3ec643e4ea4871dfe515c559d28bc0d23b6d9d6b22daf206f1de9aff13e51591323199; addinfo=%7B%22chkadmin%22%3A1%2C%22chkarticle%22%3A1%2C%22levelname%22%3A%22%5Cu7ba1%5Cu7406%5Cu5458%22%2C%22userid%22%3A%221%22%2C%22useralias%22%3A%22admin%22%7D; Hm_lvt_12d9f8f1740b76bb88c6691ea1672d8b=1589265192,1589341266,1589717172,1589769747; timezone=8`}).Debug(dbug).Do()
			assert.NoError(t, err)
			//io.Copy(os.Stdout, &buf)
			assert.Equal(t, bytes.Index(buf.Bytes(), []byte("NOVERB")), -1)
			return err
		}(),
		// formdata
		func() error {
			buf.Reset()
			key := "testFormdata"
			val := "testFormdataValue"
			err := New().POST(ts.URL).Debug(dbug).SetForm(core.H{key: val}).Do()
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(key)), -1, core.BytesToString(buf.Bytes()))
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(val)), -1)
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(rspVal)), -1)
			return err
		}(),
		// object json
		func() error {
			buf.Reset()
			key := "testkeyjson"
			val := "testvaluejson"
			err := New().POST(ts.URL).SetJSON(core.H{key: val}).Debug(dbug).Do()

			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(key)), -1, core.BytesToString(buf.Bytes()))
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(val)), -1)
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(rspVal)), -1)
			return err
		}(),
		// array json
		func() error {
			buf.Reset()
			key := "testkeyjson"
			val := "testvaluejson"
			err := New().POST(ts.URL).SetJSON(core.A{key, val}).Debug(dbug).Do()

			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(key)), -1, core.BytesToString(buf.Bytes()))
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(val)), -1)
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(rspVal)), -1)
			return err
		}(),
		// body
		func() error {
			buf.Reset()
			key := "testFormdata"
			err := New().POST(ts.URL).Debug(dbug).SetBody(key).Do()
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(key)), -1, core.BytesToString(buf.Bytes()))
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(rspVal)), -1)
			return err
		}(),
		// yaml
		func() error {
			buf.Reset()
			key := "testkeyyaml"
			val := "testvalueyaml"
			err := New().POST(ts.URL).Debug(dbug).SetYAML(core.H{key: val}).Do()
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(key)), -1, core.BytesToString(buf.Bytes()))
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(val)), -1)
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(rspVal)), -1)
			return err
		}(),
		// xml
		func() error {
			val := "testXMLValue"

			var d data
			d.Data = val
			err := New().POST(ts.URL).Debug(dbug).SetXML(d).Do()
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(val)), -1, core.BytesToString(buf.Bytes()))
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(rspVal)), -1)
			buf.Reset()
			return err
		}(),
		// x-www-form-urlencoded
		func() error {
			buf.Reset()
			key := "testwwwform"
			val := "testwwwformvalue"
			err := New().POST(ts.URL).Debug(dbug).SetWWWForm(core.H{key: val}).Do()
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(key)), -1, core.BytesToString(buf.Bytes()))
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(val)), -1)
			assert.NotEqual(t, bytes.Index(buf.Bytes(), []byte(rspVal)), -1)
			return err
		}(),
	} {
		assert.NoError(t, err, fmt.Sprintf("test index :%d", index))
		if err != nil {
			break
		}

	}
}

func TestDebug(t *testing.T) {
	buf := &bytes.Buffer{}

	router := setupDebug(t)
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

	color.NoColor = false
	test := []func() debug.Apply{
		// 测试颜色
		func() debug.Apply {
			return debug.Func(func(o *debug.Options) {
				buf.Reset()
				o.Debug = true
				o.Color = true
				o.Write = buf
			})
		},

		// 测试打开日志输出
		func() debug.Apply {
			return debug.Func(func(o *debug.Options) {
				//t.Logf("--->1.debug.Options address = %p\n", o)
				o.Debug = true
			})
		},

		// 测试修改输出源
		func() debug.Apply {
			return debug.Func(func(o *debug.Options) {
				//t.Logf("--->2.debug.Options address = %p\n", o)
				buf.Reset()
				o.Debug = true
				o.Write = buf
			})
		},

		// 测试环境变量
		func() debug.Apply {
			return debug.Func(func(o *debug.Options) {
				buf.Reset()
				if len(os.Getenv("IOS_DEBUG")) > 0 {
					o.Debug = true
				}
				o.Write = buf
			})
		},

		// 没有颜色输出
		debug.NoColor,
	}

	s := ""
	os.Setenv("IOS_DEBUG", "true")
	for k, v := range test {
		s = ""
		err := GET(ts.URL).
			Debug(v()).
			SetBody(fmt.Sprintf("%d test debug.", k)).
			BindBody(&s).
			Do()
		assert.NoError(t, err)

		if k != 0 {
			assert.NotEqual(t, buf.Len(), 0)
		}

		assert.Equal(t, fmt.Sprintf("%d test debug.", k), s)
	}

	err := GET(ts.URL).Debug(true).SetBody("true test debug").BindBody(&s).Do()

	assert.NoError(t, err)
	assert.Equal(t, s, "true test debug")

	//d := myDup{}
	err = GET(ts.URL).Debug(false).SetBody("false test debug").BindBody(&s).Do()
	//d.reset()

	//assert.Equal(t, false, d.empty())
	assert.NoError(t, err)
	assert.Equal(t, s, "false test debug")
}

//
func Test_Debug_Apply(t *testing.T) {

	router := setupDebug(t)
	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))

	fd, err := os.Create("../testdata/debug_apply.1.txt")
	if err != nil {
		panic(err.Error())
	}
	defer fd.Close()

	test := []debug.Apply{
		debug.ToWriter(fd, true),
		debug.ToFile("../testdata/debug_apply.2.txt", true),
	}

	s := ""
	for _, debugApply := range test {

		err := GET(ts.URL).Debug(debugApply).SetBody("true test debug").BindBody(&s).Do()
		assert.NoError(t, err)
		assert.Equal(t, s, "true test debug")
	}

	for _, file := range []string{
		"../testdata/debug_apply.1.txt",
		"../testdata/debug_apply.2.txt",
	} {
		all, err := ioutil.ReadFile(file)
		assert.NoError(t, err)
		assert.NotEqual(t, bytes.Index(all, []byte("true test debug")), -1)
		os.Remove(file)
	}
}
