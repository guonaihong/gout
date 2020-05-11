package dataflow

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

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
