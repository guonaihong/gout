package export

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/guonaihong/gout/core"
	"github.com/stretchr/testify/assert"
)

func Test_Export_Curl_newTemplate(t *testing.T) {
	ts := []*template.Template{newTemplate(false), newTemplate(true)}

	for _, v := range ts {
		c := curl{
			Method: "GET",
			Header: []string{"appkey:appkeyval", "sid:hello"},
			Data:   "good",
		}

		err := v.Execute(os.Stdout, c)

		assert.NoError(t, err)
	}
}

func Test_Export_Curl_GenCurl(t *testing.T) {
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/test/path?q1=v1", strings.NewReader("wo shi body"))
	req.Header.Add("h1", "v1")
	req.Header.Add("h1", "v2")

	assert.NoError(t, err)

	assert.NoError(t, GenCurl(req, true, os.Stdout))
}

func Test_Export_isExists(t *testing.T) {
	tst := []string{"../testdata/voice.pcm", "../testdata/raw-http-post-json.txt"}

	for _, fName := range tst {
		assert.Equal(t, isExists(fName), true, fName)
	}

	notExists := []string{"a.txt", "b.txt"}
	for _, fName := range notExists {
		assert.Equal(t, isExists(fName), false, fName)
	}
}

func Test_Export_formdata(t *testing.T) {
	type need struct {
		data []interface{}
		need []string
	}

	tst := []need{
		{[]interface{}{"text", "good", "voice", core.FormFile("../testdata/voice.pcm"), "mode", "a"}, []string{`"text=good"`, `"voice=@./voice"`, `"mode=a"`}},
	}

	defer func() {
		os.Remove("./voice")
		for i := 0; i < 10; i++ {
			os.Remove(fmt.Sprintf("./voice.%d", i))
		}
	}()
	genFormdata := func(data []interface{}) (*http.Request, error) {
		b := &bytes.Buffer{}
		w := multipart.NewWriter(b)

		for i := 0; i < len(data); i += 2 {
			var key string
			switch v := data[i].(type) {
			case string:
				key = v
			default:
				return nil, errors.New("fail")
			}

			switch v := data[i+1].(type) {
			case string:
				part, err := w.CreateFormField(key)
				assert.NoError(t, err)
				_, err = part.Write([]byte(v))
				assert.NoError(t, err)
				if err != nil {
					return nil, err
				}
			case core.FormFile:
				part, err := w.CreateFormFile(key, key)
				if err != nil {
					return nil, err
				}

				all, err := ioutil.ReadFile(string(v))
				if err != nil {
					return nil, err
				}
				_, err = part.Write(all)
				assert.NoError(t, err)

			default:
				return nil, errors.New("fail")
			}
		}

		w.Close()
		req, err := http.NewRequest("GET", "www.qq.com", b)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Content-Type", w.FormDataContentType())
		return req, nil

	}

	for _, tv := range tst {
		req, err := genFormdata(tv.data)
		assert.NoError(t, err)
		if err != nil {
			return
		}
		c := curl{}
		_ = c.formData(req)
		assert.NoError(t, err)
		assert.Equal(t, tv.need, c.FormData)
	}
}

func Test_Export_getFileName(t *testing.T) {
	type need struct {
		need string
		got  string
	}

	tst := []need{
		{"../testdata/voice.pcm", "../testdata/voice.pcm.0"},
		{"../testdata/raw-http-post-json.txt", "../testdata/raw-http-post-json.txt.0"},
	}

	for _, v := range tst {
		assert.Equal(t, getFileName(v.need), v.got)
	}
}
