package encode

import (
	"bytes"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testXml struct {
	I int     `xml:"i"`
	F float64 `xml:"f"`
	S string  `xml:"s"`
}

func TestNewXMLEncode(t *testing.T) {
	x := NewXMLEncode(nil)
	assert.Nil(t, x)
}

func TestXMLEncode_Encode(t *testing.T) {
	need := testXml{
		I: 100,
		F: 3.14,
		S: "test encode xml",
	}

	out := bytes.Buffer{}

	data := []interface{}{need, &need}
	for _, v := range data {
		x := NewXMLEncode(v)
		out.Reset()

		x.Encode(&out)

		got := testXml{}

		err := xml.Unmarshal(out.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, got, need)
	}
}
