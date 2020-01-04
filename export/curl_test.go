package export

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"strings"
	"testing"
	"text/template"
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

	GenCurl(req, true, os.Stdout)
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
