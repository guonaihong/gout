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
