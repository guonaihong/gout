package export

import (
	"github.com/stretchr/testify/assert"
	"os"
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
