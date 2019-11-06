package color

import (
	"github.com/guonaihong/gout/core"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

type testFormat struct {
	r         io.Reader
	openColor bool
	bodyType  BodyType
}

func Test_ColorCore_NewFormat_Nil(t *testing.T) {
	data := []testFormat{
		{nil, false, JsonType},
		{nil, true, TxtType},
		{&core.ReadCloseFail{}, true, JsonType},
		{strings.NewReader("xxx"), true, JsonType},
	}

	for _, d := range data {
		assert.Nil(t, NewFormatEncoder(d.r, d.openColor, d.bodyType))
	}
}

func Test_ColorCore_Read(t *testing.T) {
}
