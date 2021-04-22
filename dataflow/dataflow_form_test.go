package dataflow

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/guonaihong/gout/core"
	"github.com/stretchr/testify/assert"
)

func TestSetFormMap(t *testing.T) {
	reqTestForm := testForm2{
		Mode:      "A",
		Text:      "good morning",
		ReqVoice2: []byte("pcm pcm"),
		Int:       1,
		Uint:      2,
		Float32:   1.12,
		Float64:   3.14,
	}

	var err error
	reqTestForm.ReqVoice, err = ioutil.ReadFile("../testdata/voice.pcm")
	assert.NoError(t, err)

	router := setupForm2(t, reqTestForm)

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	g := New(nil)
	code := 0
	err = g.POST(ts.URL + "/test.form").
		SetForm(core.H{"mode": "A",
			"text":    "good morning",
			"voice":   core.FormFile("../testdata/voice.pcm"),
			"voice2":  core.FormMem("pcm pcm"),
			"int":     1,
			"uint":    2,
			"float32": 1.12,
			"float64": 3.14,
		}).
		Code(&code).
		Do()

	assert.NoError(t, err)
}

func TestSetFormStruct(t *testing.T) {
	reqTestForm := testForm{
		Mode:    "A",
		Text:    "good morning",
		Int:     1,
		Int8:    2,
		Int16:   3,
		Int32:   4,
		Int64:   5,
		Uint:    6,
		Uint8:   7,
		Uint16:  8,
		Uint32:  9,
		Uint64:  10,
		Float32: 1.23,
		Float64: 3.14,
		//Voice: []byte("pcm data"),
	}

	router := setupForm(t, reqTestForm)

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	code := 0
	for _, v := range [][]interface{}{
		[]interface{}{reqTestForm}, //SetForm接口 传递一个结构体
		[]interface{}{ //SetForm接口 传递二个结构体
			testForm{
				Mode:  "A",
				Text:  "good morning",
				Int:   1,
				Int8:  2,
				Int16: 3,
				Int32: 4,
			},
			testForm{
				Int64:   5,
				Uint:    6,
				Uint8:   7,
				Uint16:  8,
				Uint32:  9,
				Uint64:  10,
				Float32: 1.23,
				Float64: 3.14,
			},
		},
	} {
		err := New().POST(ts.URL + "/test.form").
			//Debug(true).
			SetForm(v...).
			Code(&code).
			Do()

		assert.NoError(t, err)
		assert.Equal(t, code, 200)

	}
}

func TestSetForm_NoAutoContentType(t *testing.T) {
	needContentType := "x-www-form-urlencoded; charset=UTF-8"
	contentTypeDefault := "application/x-www-form-urlencoded"

	router := func() *gin.Engine {
		router := gin.New()
		router.POST("/test.form", func(c *gin.Context) {

			type testFormHeader struct {
				ContentType string `header:"content-Type"`
				NeedValue   string `header:"needValue"`
			}

			var header testFormHeader
			c.ShouldBindHeader(&header)
			assert.Equal(t, header.ContentType, header.NeedValue)
			if header.ContentType == header.NeedValue {
				c.String(200, "ok")
			} else {
				c.String(500, "fail")
			}
		})
		return router
	}()

	ts := httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	defer ts.Close()

	// 1.告知gout,不加默认的content-type, 然后使用SetHeader 设置的contentType
	code := 0
	err := New().POST(ts.URL + "/test.form").
		NoAutoContentType().
		SetHeader(core.A{"content-type", needContentType, "needValue", needContentType}).
		SetWWWForm(
			core.A{
				"elapsed_time", 0,
				"user_choice_index", -1,
				"ts", 1608191572029,
			},
		).Code(&code).Do()

	assert.NoError(t, err)
	assert.Equal(t, code, 200)

	// 2.告知gout,不加默认的content-type, 然后不设置新的contentType,该字段应该为空
	err = New().POST(ts.URL + "/test.form").
		NoAutoContentType().
		//SetHeader(core.A{"content-type", needContentType}).
		SetWWWForm(
			core.A{
				"elapsed_time", 0,
				"user_choice_index", -1,
				"ts", 1608191572029,
			},
		).Code(&code).Do()

	assert.NoError(t, err)
	assert.Equal(t, code, 200)

	// 3.默认形为
	code = 0
	err = New().POST(ts.URL + "/test.form").
		SetHeader(core.A{"needValue", contentTypeDefault}).
		SetWWWForm(
			core.A{
				"elapsed_time", 0,
				"user_choice_index", -1,
				"ts", 1608191572029,
			},
		).Code(&code).Do()

	assert.NoError(t, err)
	assert.Equal(t, code, 200)

	// 4.直接设置contentType
	// 这里会有两个content-Type, 有的服务器可能会报错,尽量避免这种情况发生.
	// 如果要自定义content-Type, 先使用NoAutoContentType,取消自动加的content-Type,
	// 然后再通过SetHeader添加,可看第1种测试情况
	code = 0
	err = New().POST(ts.URL + "/test.form").
		SetHeader(core.A{"needValue", needContentType, "content-type", needContentType}).
		SetWWWForm(
			core.A{
				"elapsed_time", 0,
				"user_choice_index", -1,
				"ts", 1608191572029,
			},
		).Code(&code).Do()

	assert.NoError(t, err)
	assert.Equal(t, code, 200)
}
