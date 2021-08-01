package color

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/guonaihong/gout/core"
	"github.com/stretchr/testify/assert"
)

type testFormat struct {
	r         io.Reader
	openColor bool
	bodyType  BodyType
}

func Test_ColorCore_NewFormat_Nil(t *testing.T) {
	data := []testFormat{
		{nil, false, JSONType},
		{nil, true, TxtType},
		{&core.ReadCloseFail{}, true, JSONType},
		{strings.NewReader("xxx"), true, JSONType},
	}

	for _, d := range data {
		assert.Nil(t, NewFormatEncoder(d.r, d.openColor, d.bodyType))
	}
}

func Test_ColorCore_Read(t *testing.T) {
	j := `{
    "array": [
        "foo",
        "bar",
        "baz"
    ],
    "bool": false,
    "null": null,
    "num": 100,
    "obj": {
        "a": 1,
        "b": 2
    },
    "str": "foo"
}`

	NoColor = false
	f := NewFormatEncoder(strings.NewReader(j), true /*open color*/, JSONType)
	var out bytes.Buffer

	_, err := io.Copy(&out, f)
	assert.NoError(t, err)

	all, err := ioutil.ReadFile("./testdata/color.data")
	if err != nil {
		return
	}

	assert.Equal(t, out.Bytes(), all)
	//ioutil.WriteFile("/tmp/color.data", out.Bytes(), 0644)
}
